package types

import (
	"time"
)

// Item represents the structure of each object in the JSON array
type Item struct {
	Name        string    `json:"name"`
	ProgramURL  string    `json:"program_url"`
	URL         string    `json:"URL"`
	Count       int       `json:"count"`
	Change      int       `json:"change"`
	IsNew       bool      `json:"is_new"`
	Platform    string    `json:"platform"`
	Bounty      bool      `json:"bounty"`
	LastUpdated time.Time `json:"last_updated"`
}
