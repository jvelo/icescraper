package scrape

import (
	"net/http"
	"net/url"
	"time"
)

// ServerStatus describes the status of a server.
type ServerStatus string

const (
	StatusUnknown ServerStatus = "unknown"
	StatusUp      ServerStatus = "up"
	StatusDown    ServerStatus = "down"
)

type Target struct {
	url                url.URL
	status             ServerStatus
	lastError          error
	lastScrape         time.Time
	lastScrapeDuration time.Duration
}

type Scraper struct {
	targets  []Target
	shutdown chan struct{}
	client   *http.Client
}

func NewScrapper(urls []url.URL) *Scraper {
	var targets = make([]Target, len(urls))
	for i, u := range urls {
		targets[i] = Target{
			url: u,
		}
	}
	s := &Scraper{
		targets: targets,
		client:  &http.Client{},
	}
	return s
}

func (s *Scraper) Start() {

}
