package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Oik17/FamPay-BE/database"
	"github.com/Oik17/FamPay-BE/services"
	"github.com/Oik17/FamPay-BE/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func SearchVideos(c *fiber.Ctx) error {
	query := c.Query("id")
	db := database.DB.Db

	if query == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "false",
			"data":    "Invalid query",
			"message": "Invalid query",
		})
	}

	apiKey := utils.Config("YOUTUBE_API_KEY")
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "false",
			"data":    err.Error(),
			"message": "Invalid API Key",
		})
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

	var videos []fiber.Map
	sqlQuery := `
	INSERT INTO video (id, title, description, channelTitle, publishedAt, thumbnail, url)
	SELECT $1, $2, $3, $4, $5, $6, $7;
	`

	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			video := fiber.Map{
				"channelTitle": item.Snippet.ChannelTitle,
				"title":        item.Snippet.Title,
				"description":  item.Snippet.Description,
				"videoUrl":     item.Id.VideoId,
				"thumbnail":    item.Snippet.Thumbnails.High.Url,
				"publishedAt":  item.Snippet.PublishedAt,
			}
			video["publishedAt"], err = time.Parse(time.RFC3339, video["publishedAt"].(string))
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to parse publishedAt",
					"data":    err.Error(),
					"status":  "false",
				})
			}
			videos = append(videos, video)
			urlCheck, err := services.CheckUrlInDB(video["videoUrl"].(string))
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to check",
					"data":    err.Error(),
					"status":  false,
				})
			}
			if urlCheck {
				_, err = db.ExecContext(
					c.Context(),
					sqlQuery,
					uuid.New(),
					video["channelTitle"].(string),
					video["title"].(string),
					video["description"].(string),
					video["publishedAt"],
					video["thumbnail"].(string),
					video["videoUrl"].(string),
				)
				if err != nil {
					return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
						"message": "Failed to create video",
						"data":    err.Error(),
						"status":  "false",
					})
				}
			}
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "true",
		"data":    videos,
		"message": "Successfully fetched videos",
	})
}
