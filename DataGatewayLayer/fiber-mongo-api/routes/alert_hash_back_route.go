package routes

import (
	"mongo-api/controllers" //add this

	"github.com/gofiber/fiber/v2"
)

func AlertaHashBackRoute(app *fiber.App) {
	app.Get("/alertHashBack", controllers.GetAllAlertHashBack)
	app.Get("/alertHashBack/:id", controllers.GetAlertHashByHashBack)
	app.Post("/alertHashBack", controllers.CreateAlertHashBack)
	app.Put("/alertHashBack/:id", controllers.UpdateAlertHashBack)
	app.Delete("/alertHashBack/:id", controllers.DeleteAlertHashBack)
}
