package controllers

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/mongo"
	"mongo-api/configs"
	"mongo-api/models"
	"mongo-api/responses"
	"net/http"
	"time"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

var deviceCollection *mongo.Collection = configs.GetCollection(configs.DB, "devices")
var validateDevice = validator.New()

func CreateDevice(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var device models.Devices
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&device); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DeviceResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validateDevice.Struct(&device); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DeviceResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	/*
	   { version: 111,
	     id: 1459633408,
	     timestamp: 1685407036,
	     input1: 1,
	     input2: 1,
	     output: 0,
	     alarm_battery: false,
	     alarm_power: false,
	     sensor_error: false,
	     sensors: [ { type: 0, value: 24.870001 }, { type: 1, value: 65.970001 } ] }
	*/

	newDevice := models.Devices{
		Hash:          device.Hash,
		IdEdge:        device.IdEdge,
		IdNodo:        device.IdNodo,
		Input1:        device.Input1,
		Input2:        device.Input2,
		Output:        device.Output,
		Alarm_battery: device.Alarm_battery,
		Alarm_power:   device.Alarm_power,
		Sensor_error:  device.Sensor_error,
		Temperatura:   device.Temperatura,
		Umidade:       device.Umidade,
		CreatedAt:     device.CreatedAt,
		UpdatedAt:     device.UpdatedAt,
	}
	/*
	       newDevice := models.Device{
	           Id:       primitive.NewObjectID(),
	           IdDevice:     device.IdDevice,
	           Input1:    device.Input1,
	           Input2:    device.Input2,
	           Output:    device.Output,
	           Alarm_battery:    device.Alarm_battery,
	           Alarm_power:    device.Alarm_power,
	           Sensor_error:    device.Sensor_error,
	   	Temperatura:    device.Temperatura,
	   	Umidade:       device.Umidade,
	   	CreatedAt:     device.CreatedAt,
	   	UpdatedAt:     device.UpdatedAt,
	       }
	*/
	result, err := deviceCollection.InsertOne(ctx, newDevice)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DeviceResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.DeviceResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAllDevices(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var devices []models.Devices
	defer cancel()

	results, err := deviceCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DeviceResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var singleDevice models.Devices
		if err = results.Decode(&singleDevice); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.DeviceResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		devices = append(devices, singleDevice)
	}

	return c.Status(http.StatusOK).JSON(
		responses.DeviceResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": devices}},
	)
}
