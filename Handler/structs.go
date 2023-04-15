package Handler

import (
	"time"
)

// Set startpoint for figuring out elapsed uptime
var StartTime = time.Now()

// Global local webhookarray, soon to be replaced with firebase functionality
var tempWebhooks []WebhookObject

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

type WebhookObject struct {
	URL         string `json:"url,omitempty"`
	ISO         string `json:"country,omitempty"`
	Calls       int    `json:"calls,omitempty"`
	Invocations int    `json:"invocations,omitempty"`
	ID          string `json:"webhook_id,omitempty"`
}

type webhookOut struct {
	ID      string `json:"webhook_id"`
	Country string `json:"country"`
	Calls   int    `json:"calls"`
}
