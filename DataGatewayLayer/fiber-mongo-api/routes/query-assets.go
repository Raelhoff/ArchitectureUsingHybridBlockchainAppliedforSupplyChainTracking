package routes

import (
	"mongo-api/controllers" //add this

	"github.com/gofiber/fiber/v2"
)

func QueryRoute(app *fiber.App) {
	app.Post("/query-assets", controllers.QueryAssets)
	app.Post("/create-assets", controllers.CreateAsset)
}
