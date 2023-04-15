package Handler

import (
	"time"
)

// Set startpoint for figuring out elapsed uptime
var StartTime = time.Now()

type Country struct {
	Name       string    `json:"name"`
	ISO        string    `json:"isocode"`
	Year       []int     `json:"year"`
	Percentage []float64 `json:"percentage"`
	Borders    []string  `json:"borders"`
}

type CountryOut struct {
	Name       string  `json:"name"`
	ISO        string  `json:"isoCode"`
	Year       int     `json:"year,omitempty"`
	Percentage float64 `json:"percentage"`
}

type webhookObject struct {
	URL   string `json:"url,omitempty"`
	ISO   string `json:"country,omitempty"`
	Calls int    `json:"calls,omitempty"`
	ID    string `json:"webhook_id,omitempty"`
}
