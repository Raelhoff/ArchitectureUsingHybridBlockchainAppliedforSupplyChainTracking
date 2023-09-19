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
	"mongo-api/responses"
	"net/http"
	"time"
//	"io/ioutil"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

var alertCollection *mongo.Collection = configs.GetCollection(configs.DB, "alerts")
var validateAlert = validator.New()

func GetAllAlert(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var alertas []models.Alerta
	defer cancel()

	results, err := alertCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println( err.Error())
		return c.Status(http.StatusInternalServerError).JSON(responses.AlertaResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleAlert models.Alerta
		if err = results.Decode(&singleAlert); err != nil {
			fmt.Println( err.Error())
			return c.Status(http.StatusInternalServerError).JSON(responses.AlertaResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		alertas = append(alertas, singleAlert)
	}

	return c.Status(http.StatusOK).JSON(
		responses.AlertaResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": alertas}},
	)
}

func CreateAlert(c *fiber.Ctx) error {
	var newAlert models.Alerta
	
	fmt.Println("CreateAlert")
	
	// Log do corpo da requisição
	bodyBytes := c.Request().Body()
	fmt.Printf("Requisição recebida: %s\n", bodyBytes)

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05") // O layout define o formato desejado

	fmt.Println("Ini - Data e Hora formatadas:", formattedTime)

	//fmt.Println(c.BodyParser())
	if err := c.BodyParser(&newAlert); err != nil {
		fmt.Println("Invalid data format")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := alertCollection.InsertOne(ctx, newAlert)
	if err != nil {
		fmt.Println(err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create alert: %v", err),
		})
	}

	currentTime2 := time.Now()
	formattedTime2 := currentTime2.Format("2006-01-02 15:04:05") // O layout define o formato desejado
	fmt.Println("Fim - Data e Hora formatadas:", formattedTime2)

	return c.Status(http.StatusCreated).JSON(newAlert)
}

func GetAlertByHash(c *fiber.Ctx) error {
	alertID := c.Params("id")

	var alert models.Alerta
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := alertCollection.FindOne(ctx, bson.M{"hash": alertID}).Decode(&alert)
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

func UpdateAlert(c *fiber.Ctx) error {
	alertID := c.Params("id")

	var updatedAlert models.Alerta
	if err := c.BodyParser(&updatedAlert); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid data format",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := alertCollection.UpdateOne(
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

func DeleteAlert(c *fiber.Ctx) error {
	alertID := c.Params("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := alertCollection.DeleteOne(ctx, bson.M{"_id": alertID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to delete alert: %v", err),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Alert deleted successfully",
	})
}
