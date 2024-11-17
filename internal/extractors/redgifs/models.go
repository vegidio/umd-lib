package redgifs

import "time"

type Video struct {
	Author  string
	Url     string
	Created time.Time
}
