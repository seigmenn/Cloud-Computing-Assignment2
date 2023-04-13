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
