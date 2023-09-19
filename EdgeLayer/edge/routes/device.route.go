package routes

import (
	"fiber-mongo-api/controllers" //add this

	"github.com/gofiber/fiber/v2"
	// pb "github.com/Raelhoff/gRPC_GO/protofiles"
)

func DeviceRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Post("/lora", controllers.CreateDevice)              //add this
	app.Get("/device/:deviceId", controllers.FindDeviceById) //add this
	//app.Put("/user/:userId", controllers.EditAUser) //add this
	app.Delete("/device/:deviceId", controllers.DeleteDevice) //add this

	app.Get("/devices", controllers.GetAllDevices) //add this
}
