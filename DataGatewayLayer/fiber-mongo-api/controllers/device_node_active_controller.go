package controllers

import (
	"context"
	"mongo-api/configs"
	"mongo-api/models"
	"mongo-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

var NodeActiceCollection *mongo.Collection = configs.GetCollection(configs.DB, "nodo_active")
var validateNodeActice = validator.New()

func GetAllDevicesNodeActive(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var devices []models.NodeActive
	defer cancel()

	results, err := NodeActiceCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.NodeActiveResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleNodeActive models.NodeActive
		if err = results.Decode(&singleNodeActive); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.NodeActiveResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		devices = append(devices, singleNodeActive)
	}

	return c.Status(http.StatusOK).JSON(
		responses.NodeActiveResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": devices}},
	)
}
