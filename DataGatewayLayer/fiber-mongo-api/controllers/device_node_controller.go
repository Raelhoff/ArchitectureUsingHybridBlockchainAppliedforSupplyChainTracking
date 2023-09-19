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

var deviceNodeCollection *mongo.Collection = configs.GetCollection(configs.DB, "devices_nodo")
var validateNodeEdge = validator.New()

func CreateDeviceNodo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var device models.DevicesNodo
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&device); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DeviceNodeResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validateNodeEdge.Struct(&device); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DeviceNodeResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	/*
		ID     primitive.ObjectID `bson:"_id,omitempty"`
		IdNodo int32              `json:"IdNodo" validate:"required"`
		IdEdge int32              `json:"IdEdge" validate:"required"`
		Period int                `json:"Period"`
		Date   string             `gorm:"not null;default:'1970-01-01 00:00:01'" json:"Date,omitempty"`
		Hash   string             `json:"Hash"`
	*/

	newDevice := models.DevicesNodo{
		IdEdge: device.IdEdge,
		IdNodo: device.IdNodo,
		ID:     device.ID,
		Period: device.Period,
		Date:   device.Date,
		Hash:   device.Hash,
	}

	result, err := deviceNodeCollection.InsertOne(ctx, newDevice)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DeviceNodeResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.DeviceNodeResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAllDevicesNodo(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var devices []models.DevicesNodo
	defer cancel()

	results, err := deviceNodeCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DeviceNodeResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleDeviceNodo models.DevicesNodo
		if err = results.Decode(&singleDeviceNodo); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.DeviceNodeResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		devices = append(devices, singleDeviceNodo)
	}

	return c.Status(http.StatusOK).JSON(
		responses.DeviceNodeResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": devices}},
	)
}
