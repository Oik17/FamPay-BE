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
	app.Listen(":3000")
}
