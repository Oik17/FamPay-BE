package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Oik17/FamPay-BE/database"
	"github.com/Oik17/FamPay-BE/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func FetchAndStoreVideos(query string, videos *[]fiber.Map) error {
	for i := 0; i < 2; i++ { //Value limited to 2 to not exhaust api credits. Can run indefinitely asynchronously by replacing this line with just

		apiKeysString := utils.Config("YOUTUBE_API_KEYS")
		apiKeys := strings.Split(apiKeysString, ",")

		for _, apiKey := range apiKeys {
			youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
			if err != nil {
				return err
			}
			if youtubeService == nil {
				return errors.New("youtube service is nil")
			}
			publishedAfter := time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC)
			call := youtubeService.Search.List([]string{"snippet"}).
				Q(query).
				Type("video").
				Order("date").
				PublishedAfter(publishedAfter.Format(time.RFC3339))

			response, err := call.Do()
			if err != nil {
				continue //if API Key is invalid, check with next APIKey
			}

			db := database.DB.Db
			sqlQuery := `INSERT INTO video (id, prompt, title, description, channelTitle, publishedAt, thumbnail, url) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

			for _, item := range response.Items {
				if item.Id.Kind == "youtube#video" {
					video := map[string]interface{}{
						"channelTitle": item.Snippet.ChannelTitle,
						"title":        item.Snippet.Title,
						"description":  item.Snippet.Description,
						"videoUrl":     item.Id.VideoId,
						"thumbnail":    item.Snippet.Thumbnails.High.Url,
						"publishedAt":  item.Snippet.PublishedAt,
					}

					video["publishedAt"], err = time.Parse(time.RFC3339, video["publishedAt"].(string))
					if err != nil {
						fmt.Printf("Failed to parse publishedAt for video %s: %v\n", item.Id.VideoId, err)
						continue
					}
					urlCheck, err := CheckUrlInDB(video["videoUrl"].(string))
					if err != nil {
						fmt.Printf("Failed to check URL in DB for video %s: %v\n", item.Id.VideoId, err)
						continue
					}

					if urlCheck {
						*videos = append(*videos, video)

						_, err = db.ExecContext(
							context.Background(),
							sqlQuery,
							uuid.New(),
							query,
							video["channelTitle"].(string),
							video["title"].(string),
							video["description"].(string),
							video["publishedAt"],
							video["thumbnail"].(string),
							video["videoUrl"].(string),
						)
						if err != nil {
							return err
						}

					}
				}
			}
			break
		}

		time.Sleep(10 * time.Second)
	}
	return nil
}
