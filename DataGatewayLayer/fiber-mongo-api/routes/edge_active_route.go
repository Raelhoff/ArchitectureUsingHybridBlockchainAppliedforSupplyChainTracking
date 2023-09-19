package routes

import (
	"mongo-api/controllers" //add this

	"github.com/gofiber/fiber/v2"
)

func EdgeActiveRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Get("/devicesEdgeActives", controllers.GetAllDevicesEdgeActive) //add this
}
