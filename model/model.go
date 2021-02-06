package model

import "github.com/jvelo/icecast-monitor/prisma/db"

type Cast = db.InnerCast
type Track = db.InnerTrack

type Record struct {
	Cast  *Cast
	Track *Track
}

func NewCast(name string, description string, url string) *Cast {
	return &db.InnerCast{
		Name:        name,
		URL:         url,
		Description: &description,
	}
}

func NewTrack(title string, listeners int) *Track {
	return &db.InnerTrack{
		Title:     title,
		Listeners: listeners,
	}
}
