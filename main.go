package main

import (
	"github.com/Oik17/FamPay-BE/controllers"
	"github.com/Oik17/FamPay-BE/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	database.Connect()
	app := fiber.New()
	app.Get("/search", controllers.SearchVideos)
	app.Get("/getVideos", controllers.GetVideos)
	app.Get("/getVideos/title", controllers.GetVideoByTitle)
	app.Get("/getVideos/id", controllers.GetVideoById)
	app.Get("/getVideos/prompt", controllers.GetVideoByPrompt)
	app.Listen(":3000")

}
