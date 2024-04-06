package main

import (
	"log"

	"github.com/Oik17/FamPay-BE/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	app.Get("/search", controllers.SearchVideos)
	app.Listen(":3000")
}
