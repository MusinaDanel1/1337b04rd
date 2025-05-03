package models

import "time"

type Session struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	LastActivity time.Time `json:"last_activity"`
}
