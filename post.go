package simpleapm

import "time"

type Post struct {
	Type      string    `json:"type"`
	Status    int       `json:"status"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}
