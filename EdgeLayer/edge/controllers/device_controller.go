package controllers

import (
	"context"
	initializers "fiber-mongo-api/configs" //add this
	"fiber-mongo-api/models"
	"fiber-mongo-api/packge"
	"fiber-mongo-api/responses"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

var validateDevice = validator.New()

func salveCreateDeviceDB(newDevice models.Devices) *gorm.DB {
	return initializers.DB.Create(&newDevice)
}

var debugTime bool = true // Defina como true para gravar tempos em arquivo de debug
var idEDGE int32 = 1234   // Defina idEDGE

func writeDebugTimingCreateDevice(timingStr string) error {
	fileName := fmt.Sprintf("/home/workspace/bitburket/simulator/logs/debug_timing_CreateDevice_%d.txt", idEDGE)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(timingStr)
	if err != nil {
		return err
	}

	return nil
}

func CreateDevice(c *fiber.Ctx) error {
	var device models.Device

	//validate the request body
	if err := c.BodyParser(&device); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DeviceResponse{Status: http.StatusBadRequest, Message: "error: BodyParser", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validateDevice.Struct(&device); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DeviceResponse{Status: http.StatusBadRequest, Message: "error: validation", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	// Salve o tempo de envio
	sendTime := time.Now()

	var conn *grpc.ClientConn = initializers.ReturnClientGRPC()
	client := packge.NewLoraTransactionClient(conn)

	loc, _ := time.LoadLocation("America/Sao_Paulo")
	t := time.Now().In(loc)

	newDevice := models.Devices{
		IdEdge:        idEDGE,
		IdNodo:        device.IdDevice,
		Input1:        device.Input1,
		Input2:        device.Input2,
		Output:        device.Output,
		Alarm_battery: device.Alarm_battery,
		Alarm_power:   device.Alarm_power,
		Sensor_error:  device.Sensor_error,
		Temperatura:   device.Sensors_obj[0].Value,
		Umidade:       device.Sensors_obj[1].Value,
		CreatedAt:     t.Format(time.RFC3339),
		UpdatedAt:     t.Format(time.RFC3339),
	}

	packge.SetDeviceStatus(int(device.IdDevice), true)

	if client != nil {
		go func() {
			tx, err := client.MakeTransaction(
				context.Background(),
				&packge.LoraRequest{
					IdEdge:       idEDGE,
					IdNodo:       device.IdDevice,
					Input1:       device.Input1,
					Input2:       device.Input2,
					Output:       device.Output,
					AlarmBattery: device.Alarm_battery,
					AlarmPower:   device.Alarm_power,
					SensorError:  device.Sensor_error,
					Sensors: []*packge.LoraRequest_Sensor{
						{Type: 1, Value: device.Sensors_obj[0].Value},
						{Type: 2, Value: device.Sensors_obj[1].Value},
					},
					LastUpdated: t.Format(time.RFC3339),
				},
			)

			// Salve o tempo de resposta
			receiveTime := time.Now()

			if err != nil {
				fmt.Println("Error, %v", err)

				// Calcule a diferença de tempo
				timeTaken := receiveTime.Sub(sendTime)

				// Converta a estrutura em uma representação de string
				timingStr := fmt.Sprintf("Confirmation:false, IdEdge: %d, IdDevice: %d, Enviado em: %s, Recebido em: %s, Tempo de resposta: %s\n",
					idEDGE, device.IdDevice, sendTime.Format("02/01/2006 15:04:05"), receiveTime.Format("02/01/2006 15:04:05"), timeTaken.String())

				// Escreva a informação de tempo no arquivo de debug
				if err := writeDebugTimingCreateDevice(timingStr); err != nil {
					fmt.Println("Erro ao gravar no arquivo de debug:", err)
				}
			} else {

				// Calcule a diferença de tempo
				timeTaken := receiveTime.Sub(sendTime)

				// Converta a estrutura em uma representação de string
				timingStr := fmt.Sprintf("Confirmation: true, IdEdge: %d, IdDevice: %d, Enviado em: %s, Recebido em: %s, Tempo de resposta: %s\n",
					idEDGE, device.IdDevice, sendTime.Format("02/01/2006 15:04:05"), receiveTime.Format("02/01/2006 15:04:05"), timeTaken.String())

				// Escreva a informação de tempo no arquivo de debug
				if err := writeDebugTimingCreateDevice(timingStr); err != nil {
					fmt.Println("Erro ao gravar no arquivo de debug:", err)
				}
				fmt.Println("Msg: ", tx.Msg)
				fmt.Println("Confirmation: ", tx.Confirmation)
			}
		}()
	} else {
		go func() {
			tx, err := client.MakeTransaction(
				context.Background(),
				&packge.LoraRequest{
					IdEdge:       idEDGE,
					IdNodo:       device.IdDevice,
					Input1:       device.Input1,
					Input2:       device.Input2,
					Output:       device.Output,
					AlarmBattery: device.Alarm_battery,
					AlarmPower:   device.Alarm_power,
					SensorError:  device.Sensor_error,
					Sensors: []*packge.LoraRequest_Sensor{
						{Type: 1, Value: device.Sensors_obj[0].Value},
						{Type: 2, Value: device.Sensors_obj[1].Value},
					},
					LastUpdated: t.Format(time.RFC3339),
				},
			)

			if err != nil {
				fmt.Println("Error, %v", err)

				result := salveCreateDeviceDB(newDevice)

				if result.Error != nil && strings.Contains(result.Error.Error(), "Duplicate entry") {
					fmt.Println("Title already exist, please use another title")
				} else if result.Error != nil {
					fmt.Println(result.Error.Error())
				}
			} else {

				fmt.Println("Msg: ", tx.Msg)
				fmt.Println("Confirmation: ", tx.Confirmation)

			}
		}()

	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": fiber.Map{"device": newDevice}})
}

func DeleteDevice(c *fiber.Ctx) error {
	deviceId := c.Params("deviceId")

	result := initializers.DB.Delete(&models.Devices{}, "inserted_id = ?", deviceId)

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that Id exists"})
	} else if result.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": result.Error})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func FindDeviceById(c *fiber.Ctx) error {
	deviceId := c.Params("deviceId")

	var device models.Devices
	result := initializers.DB.First(&device, "inserted_id = ?", deviceId)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "No note with that Id exists"})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"device": device}})
}

func GetAllDevices(c *fiber.Ctx) error {
	rows, err := initializers.DB2.Query("SELECT * FROM devices")

	if err != nil {
		return c.Status(http.StatusOK).JSON(
			responses.DeviceResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": err}},
		)
	}
	defer rows.Close()

	var devices = make([]models.Devices, 0)

	for rows.Next() {
		singleDevice := models.Devices{}
		err = rows.Scan(&singleDevice.InsertedID,
			&singleDevice.IdEdge, &singleDevice.IdNodo, &singleDevice.Input1,
			&singleDevice.Input2, &singleDevice.Output,
			&singleDevice.Alarm_battery, &singleDevice.Alarm_power,
			&singleDevice.Sensor_error, &singleDevice.Temperatura,
			&singleDevice.Umidade, &singleDevice.CreatedAt,
			&singleDevice.UpdatedAt)

		if err != nil {
			return c.Status(http.StatusOK).JSON(
				responses.DeviceResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": err}},
			)
		}

		devices = append(devices, singleDevice)
	}

	err = rows.Err()

	if err != nil {
		return c.Status(http.StatusOK).JSON(
			responses.DeviceResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": err}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.DeviceResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": devices}},
	)
}
