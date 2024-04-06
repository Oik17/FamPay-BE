package utils

import (
	"context"
	"errors"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func YoutubeService(apiKeys []string) (*youtube.Service, error) {
	for _, apiKey := range apiKeys {
		youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
		if err == nil {
			return youtubeService, nil
		}
		log.Println(err.Error())
	}

	return nil, errors.New("failed to create YouTube service with all provided API keys")
}
