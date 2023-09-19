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

var EdgeActiceCollection *mongo.Collection = configs.GetCollection(configs.DB, "edge_active")
var validateEdgeActice = validator.New()

func GetAllDevicesEdgeActive(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var devices []models.EdgeActive
	defer cancel()

	results, err := EdgeActiceCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.EdgeActiveResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleEdgeActive models.EdgeActive
		if err = results.Decode(&singleEdgeActive); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.EdgeActiveResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		devices = append(devices, singleEdgeActive)
	}

	return c.Status(http.StatusOK).JSON(
		responses.EdgeActiveResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": devices}},
	)
}
