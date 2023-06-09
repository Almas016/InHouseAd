package model

import "time"

type Website struct {
	URL        string
	AccessTime time.Duration
}
