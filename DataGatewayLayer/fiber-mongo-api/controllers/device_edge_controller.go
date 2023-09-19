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

var deviceEdgeCollection *mongo.Collection = configs.GetCollection(configs.DB, "devices_edge")
var validateDeviceEdge = validator.New()

func CreateDeviceEdge(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var device models.DevicesEdge
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&device); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DeviceEdgeResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validateDevice.Struct(&device); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DeviceEdgeResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	/*
			IdEdge int32              `json:"IdEdge" validate:"required"`
		Ip     string             `json:"Ip"`
		User   string             `json:"User"`
		Date   string             `gorm:"not null;default:'1970-01-01 00:00:01'" json:"Date,omitempty"`
		Hash   string             `json:"Hash"`
	*/

	newDevice := models.DevicesEdge{
		Hash:   device.Hash,
		IdEdge: device.IdEdge,
		Ip:     device.Ip,
		User:   device.User,
		Date:   device.Date,
	}

	result, err := deviceEdgeCollection.InsertOne(ctx, newDevice)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DeviceEdgeResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.DeviceEdgeResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAllDevicesEdge(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var devices []models.DevicesEdge
	defer cancel()

	results, err := deviceEdgeCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DeviceEdgeResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleDeviceEdge models.DevicesEdge
		if err = results.Decode(&singleDeviceEdge); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.DeviceEdgeResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		devices = append(devices, singleDeviceEdge)
	}

	return c.Status(http.StatusOK).JSON(
		responses.DeviceEdgeResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": devices}},
	)
}
