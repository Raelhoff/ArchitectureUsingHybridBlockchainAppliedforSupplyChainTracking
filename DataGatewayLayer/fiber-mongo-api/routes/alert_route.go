package routes

import (
	"mongo-api/controllers" //add this

	"github.com/gofiber/fiber/v2"
)

func AlertaRoute(app *fiber.App) {
	app.Get("/alert", controllers.GetAllAlert)
	app.Get("/alert/:id", controllers.GetAlertByHash)
	app.Post("/alert", controllers.CreateAlert)
	app.Put("/alert/:id", controllers.UpdateAlert)
	app.Delete("/alert/:id", controllers.DeleteAlert)
}
