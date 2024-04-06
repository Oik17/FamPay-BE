package models

import (
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID           uuid.UUID `json:"id"`
	Prompt       string    `json:"prompt"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ChannelTitle string    `json:"channelTitle"`
	PublishedAt  time.Time `json:"publishedAt"`
	Thumbnail    string    `json:"thumbnail"`
	URL          string    `json:"url"`
}
