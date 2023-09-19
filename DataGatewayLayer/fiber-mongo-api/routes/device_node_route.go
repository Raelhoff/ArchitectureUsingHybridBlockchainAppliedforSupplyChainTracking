package routes

import (
	"mongo-api/controllers" //add this

	"github.com/gofiber/fiber/v2"
)

func DeviceNodeRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Post("/deviceNode", controllers.CreateDeviceNodo) //add this
	// app.Get("/user/:userId", controllers.GetDevice) //add this
	//app.Put("/user/:userId", controllers.EditAUser) //add this
	//app.Delete("/user/:userId", controllers.DeleteAUser) //add this

	app.Get("/deviceNode", controllers.GetAllDevicesNodo) //add this
}
