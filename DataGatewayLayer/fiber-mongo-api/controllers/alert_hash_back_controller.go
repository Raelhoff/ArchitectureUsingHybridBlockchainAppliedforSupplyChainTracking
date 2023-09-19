package controllers

import (
	"context"
	"fmt"
	"mongo-api/configs"
	"mongo-api/models"
	"mongo-api/packge"
	"mongo-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"crypto/sha512"
	"encoding/hex"
)

var alertHashBackCollection *mongo.Collection = configs.GetCollection(configs.DB, "alerts")
var alertHashFifoBackCollection *mongo.Collection = configs.GetCollection(configs.DB, "alertsHashFiFo")
var validateAlertHashBack = validator.New()

func InteractWithSmartContractBack(newAlert models.Alerta) {
	go func() {

		contractSmartHash := packge.GetContractSmartHash()

		var ok = ConvertAlertaToAssetHash(newAlert)
		// Execute the contract method (adjust parameters as needed)
		_, err := contractSmartHash.SubmitTransaction("CreateAsset",
			ok.ID, ok.IdDevice, ok.Type,
			ok.ProductionDate, ok.ProductionLocation, ok.Description)
		if err != nil {
			//	return err
		}

		// Apenas para fins de exemplo, aqui você pode atualizar o alerta no banco de dados para marcar como processado
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := alertHashBackCollection.UpdateOne(
			ctx,
			bson.M{"hash": newAlert.Hash},
			bson.M{"$set": bson.M{"processed": true}},
		)
		if err != nil {
			fmt.Printf("Failed to update alert in database: %v", err)
		}
		fmt.Println(result.ModifiedCount)
	}()

}

func ConvertAlertaToAssetHashBack(alerta models.AssetHashBack) models.AssetHash {
	assetHash := models.AssetHash{
		ID:                 alerta.Hash,
		IdDevice:           fmt.Sprintf("%d", alerta.IdEdge),
		Type:               fmt.Sprintf("%d", alerta.Type),
		ProductionDate:     alerta.Date,
		ProductionLocation: "São José",
		Description:        alerta.Description,
	}
	return assetHash
}

func CreateAlertHashBack(c *fiber.Ctx) error {
	var newAlert models.AssetHashBack

	fmt.Println("CreateAlertHashBack")

	// Log do corpo da requisição
	bodyBytes := c.Request().Body()
	fmt.Printf("Requisição recebida: %s\n", bodyBytes)

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05") // O layout define o formato desejado

	fmt.Println("Ini - Data e Hora formatadas:", formattedTime)

	if err := c.BodyParser(&newAlert); err != nil {
		fmt.Println("Invalid data format")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var hashdevice = ConvertAlertaToAssetHashBack(newAlert)

	dataOut := EncodeToBytes(hashdevice)
	//fmt.Println([]byte(dataOut))
	sha_512 := sha512.New()
	// sha from a byte array
	sha_512.Write([]byte(dataOut))
	fmt.Printf("sha512: %x\n", sha_512.Sum(nil))
	newAlert.Hash = hex.EncodeToString(sha_512.Sum(nil))

	_, err := alertHashBackCollection.InsertOne(ctx, newAlert)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create alert: %v", err),
		})
	}

	// Interact with the SmartHash contract
	//InteractWithSmartContractBack(newAlert)

	// Insert only the hash into the alertHashFifoBackCollection
	hashDocument := bson.M{"hash": newAlert.Hash}
	_, err = alertHashFifoBackCollection.InsertOne(ctx, hashDocument)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create alert hash in FIFO collection: %v", err),
		})
	}

	currentTime2 := time.Now()
	formattedTime2 := currentTime2.Format("2006-01-02 15:04:05") // O layout define o formato desejado
	fmt.Println("Fim - Data e Hora formatadas:", formattedTime2)

	return c.Status(http.StatusCreated).JSON(newAlert)
}

func GetAllAlertHashBack(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var alertas []models.Alerta
	defer cancel()

	results, err := alertHashBackCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.AlertaResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleAlert models.Alerta
		if err = results.Decode(&singleAlert); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.AlertaResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		alertas = append(alertas, singleAlert)
	}

	return c.Status(http.StatusOK).JSON(
		responses.AlertaResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": alertas}},
	)
}

func GetAlertHashByHashBack(c *fiber.Ctx) error {
	alertID := c.Params("id")

	var alert models.Alerta
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := alertHashBackCollection.FindOne(ctx, bson.M{"hash": alertID}).Decode(&alert)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Alert not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get alert: %v", err),
		})
	}

	return c.JSON(alert)
}

func UpdateAlertHashBack(c *fiber.Ctx) error {
	alertID := c.Params("id")

	var updatedAlert models.Alerta
	if err := c.BodyParser(&updatedAlert); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := alertHashBackCollection.UpdateOne(
		ctx,
		bson.M{"_id": alertID},
		bson.M{"$set": updatedAlert},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to update alert: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Alert updated successfully",
	})
}

func DeleteAlertHashBack(c *fiber.Ctx) error {
	alertID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := alertHashBackCollection.DeleteOne(ctx, bson.M{"_id": alertID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to delete alert: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Alert deleted successfully",
	})
}
