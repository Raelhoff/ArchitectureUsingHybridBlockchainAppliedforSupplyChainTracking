package controllers

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/mongo"
	"mongo-api/configs"
	"mongo-api/models"
	"mongo-api/packge"
	"mongo-api/responses"
	"net/http"
	"time"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"bytes"
	"crypto/sha512"
	"encoding/gob"
	"encoding/hex"
	"log"
)

var alertHashCollection *mongo.Collection = configs.GetCollection(configs.DB, "alerts")
var validateAlertHash = validator.New()

func CreateAlertHash(c *fiber.Ctx) error {
	var newAlert models.Alerta
	if err := c.BodyParser(&newAlert); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var hashdevice = ConvertAlertaToAssetHash(newAlert)

	dataOut := EncodeToBytes(hashdevice)
	//fmt.Println([]byte(dataOut))
	sha_512 := sha512.New()
	// sha from a byte array
	sha_512.Write([]byte(dataOut))
	fmt.Printf("sha512: %x\n", sha_512.Sum(nil))
	newAlert.Hash = hex.EncodeToString(sha_512.Sum(nil))
	hashdevice.ID = newAlert.Hash

	_, err := alertHashCollection.InsertOne(ctx, newAlert)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create alert: %v", err),
		})
	}

	// Interact with the SmartHash contract
	err = interactWithSmartContract(hashdevice)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to interact with smart contract: %v", err),
		})
	}

	return c.Status(http.StatusCreated).JSON(newAlert)
}

func interactWithSmartContract(newAlert models.AssetHash) error {

	// Obter a instância do contrato
	contractSmartHash := packge.GetContractSmartHash()

	// Execute the contract method (adjust parameters as needed)
	_, err := contractSmartHash.SubmitTransaction("CreateAsset",
		newAlert.ID, newAlert.IdDevice, newAlert.Type,
		newAlert.ProductionDate, newAlert.ProductionLocation, newAlert.Description)
	if err != nil {
		return err
	}

	return nil
}

func ConvertAlertaToAssetHash(alerta models.Alerta) models.AssetHash {
	assetHash := models.AssetHash{
		ID:                 alerta.Hash,
		IdDevice:           fmt.Sprintf("%d", alerta.IDLocal),
		Type:               fmt.Sprintf("%d", alerta.Type),
		ProductionDate:     alerta.Date,
		ProductionLocation: "São José",
		Description:        alerta.Description,
	}
	return assetHash
}

func EncodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes()
}

func GetAllAlertHash(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var alertas []models.Alerta
	defer cancel()

	results, err := alertHashCollection.Find(ctx, bson.M{})
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

func GetAlertHashByHash(c *fiber.Ctx) error {
	alertID := c.Params("id")

	var alert models.Alerta
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := alertHashCollection.FindOne(ctx, bson.M{"hash": alertID}).Decode(&alert)
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

func UpdateAlertHash(c *fiber.Ctx) error {
	alertID := c.Params("id")

	var updatedAlert models.Alerta
	if err := c.BodyParser(&updatedAlert); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := alertHashCollection.UpdateOne(
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

func DeleteAlertHash(c *fiber.Ctx) error {
	alertID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := alertHashCollection.DeleteOne(ctx, bson.M{"_id": alertID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to delete alert: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Alert deleted successfully",
	})
}
