package routes

import (
	"mongo-api/controllers" //add this

	"github.com/gofiber/fiber/v2"
)

func AlertaHashRoute(app *fiber.App) {
	app.Get("/alertHash", controllers.GetAllAlertHash)
	app.Get("/alertHash/:id", controllers.GetAlertHashByHash)
	app.Post("/alertHash", controllers.CreateAlertHash)
	app.Put("/alertHash/:id", controllers.UpdateAlertHash)
	app.Delete("/alertHash/:id", controllers.DeleteAlertHash)
}
