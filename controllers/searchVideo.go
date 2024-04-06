package controllers

import (
	"net/http"

	"github.com/Oik17/FamPay-BE/services"

	"github.com/Oik17/FamPay-BE/database"
	"github.com/Oik17/FamPay-BE/models"
	"github.com/gofiber/fiber/v2"
)

func SearchVideos(c *fiber.Ctx) error {
	query := c.Query("id")
	if query == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "false",
			"data":    "Invalid query",
			"message": "Invalid query",
		})
	}

	var videos []fiber.Map
	go services.FetchAndStoreVideos(query, &videos) //Keep running the function asynchronously even after the value is returned.

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "true",
		"data":    videos,
		"message": "Successfully fetched videos",
	})
}

func GetVideos(c *fiber.Ctx) error {
	db := database.DB.Db
	var videos []models.Video
	query := `SELECT * FROM video`

	err := db.Select(&videos, query)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch videos",
			"data":    err.Error(),
			"status":  "false",
		})
	}
	return c.Status(http.StatusAccepted).JSON(videos)
}

func GetVideoByPrompt(c *fiber.Ctx) error {
	id := c.Query("id")
	db := database.DB.Db
	var videos []models.Video
	query := `SELECT * FROM video where prompt=$1`
	err := db.Select(&videos, query, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch video",
			"data":    err.Error(),
			"status":  "false",
		})
	}
	return c.Status(http.StatusAccepted).JSON(videos)
}
func GetVideoById(c *fiber.Ctx) error {
	id := c.Query("id")
	db := database.DB.Db
	var videos []models.Video

	query := `SELECT * FROM video WHERE id=$1`
	err := db.Select(&videos, query, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch video",
			"data":    err.Error(),
			"status":  "false",
		})
	}

	if len(videos) == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Video not found",
			"status":  "false",
		})
	}

	return c.Status(http.StatusOK).JSON(videos)
}

func GetVideoByTitle(c *fiber.Ctx) error {
	id := c.Query("id")
	db := database.DB.Db
	var videos []models.Video
	query := `SELECT * FROM video where title=$1`
	err := db.Select(&videos, query, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch video",
			"data":    err.Error(),
			"status":  "false",
		})
	}
	return c.Status(http.StatusAccepted).JSON(videos)
}
