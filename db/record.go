package db

import (
	"github.com/jvelo/icecast-monitor/models"
	"github.com/volatiletech/null/v8"
)

type Cast = models.Cast
type Track = models.Track

type Record struct {
	Cast  *Cast
	Track *Track
}

func NewCast(name string, description string, url string) *Cast {
	return &Cast{
		Name:        name,
		URL:         url,
		Description: null.StringFrom(description),
	}
}

func NewTrack(title string, listeners int) *Track {
	return &Track{
		Title:     title,
		Listeners: listeners,
	}
}
