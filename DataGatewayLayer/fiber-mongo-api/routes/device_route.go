package routes

import (
    "mongo-api/controllers" //add this
    "github.com/gofiber/fiber/v2"
)

func DeviceRoute(app *fiber.App) {
    //All routes related to users comes here
    app.Post("/device", controllers.CreateDevice) //add this
   // app.Get("/user/:userId", controllers.GetDevice) //add this
    //app.Put("/user/:userId", controllers.EditAUser) //add this
    //app.Delete("/user/:userId", controllers.DeleteAUser) //add this
 
    app.Get("/device", controllers.GetAllDevices) //add this
}
