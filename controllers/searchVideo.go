package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
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

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	publishedAfter := time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC)

	call := youtubeService.Search.List([]string{"snippet"}).
		Q(query).
		Type("video").
		Order("date").
		PublishedAfter(publishedAfter.Format(time.RFC3339))

	response, err := call.Do()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "false",
			"data":    err.Error(),
			"message": "Failed to fetch videos",
		})
	}

	var result string
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			result += fmt.Sprintf("Title: %s, Video ID: %s\n", item.Snippet.Title, item.Id.VideoId)
		}
	}

	return c.SendString(result)
}
