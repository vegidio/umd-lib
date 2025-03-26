package imaglr

import "time"

type Post struct {
	Id        string
	Author    string
	Type      string
	Image     string
	Video     string
	Timestamp time.Time
}
