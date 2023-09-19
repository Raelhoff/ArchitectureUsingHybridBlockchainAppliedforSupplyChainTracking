package routes

import (
	"github.com/gofiber/fiber/v2"
	"protofiles/server/controllers" //add this
)

func DeviceRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Post("/lora", controllers.CreateDevice) //add this
	// app.Get("/user/:userId", controllers.GetDevice) //add this
	//app.Put("/user/:userId", controllers.EditAUser) //add this
	//app.Delete("/user/:userId", controllers.DeleteAUser) //add this

	app.Get("/allDeviceLora", controllers.GetAllDevices) //add this
}
