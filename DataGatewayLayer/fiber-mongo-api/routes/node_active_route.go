package routes

import (
	"mongo-api/controllers" //add this

	"github.com/gofiber/fiber/v2"
)

func NodeActiveRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Get("/devicesNodeActives", controllers.GetAllDevicesNodeActive) //add this
}
