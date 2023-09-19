package utils

import (
	"context"
	"fiber-mongo-api/models"
	"fmt"
	"log"

	initializers "fiber-mongo-api/configs" //add this

	pb "github.com/Raelhoff/gRPC_GO/protofiles"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func test() {

	log.Println("🚀 Connected Successfully to the utils")
}

func SendDevices(c *fiber.Ctx) error {
	var device models.Device
	var conn *grpc.ClientConn = initializers.ReturnClientGRPC()
	client := pb.NewLoraTransactionClient(conn)

	if client != nil {
		fmt.Println("Client ok!!")

		// Lets invoke the remote function from client on the server
		tx, err := client.MakeTransaction(
			context.Background(),
			&pb.LoraRequest{
				IdDevice:     1234,
				Id:           device.IdDevice,
				Input1:       device.Input1,
				Input2:       device.Input2,
				Output:       device.Output,
				AlarmBattery: device.Alarm_battery,
				AlarmPower:   device.Alarm_power,
				SensorError:  device.Sensor_error,
				Sensors: []*pb.LoraRequest_Sensor{
					{Type: (device.Sensors_obj[0].Type + 1), Value: device.Sensors_obj[0].Value},
					{Type: (device.Sensors_obj[1].Type + 1), Value: device.Sensors_obj[1].Value},
				},
				LastUpdated: timestamppb.Now(),
			})

		if err != nil {
			fmt.Println("Error, %v", err)
		}

		fmt.Println("Msg: ", tx.Msg)
		fmt.Println("Confirmation: ", tx.Confirmation)

		return nil
	}
	return nil
}
