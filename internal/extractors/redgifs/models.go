package redgifs

import "github.com/vegidio/umd-lib/internal"

type Auth struct {
	Token   string `json:"token"`
	Session string `json:"session"`
}

type Video struct {
	Gif Gif `json:"gif"`
}

type Gif struct {
	Id       string             `json:"id"`
	Username string             `json:"userName"`
	Duration float64            `json:"duration"`
	Url      Url                `json:"urls"`
	Created  internal.EpochTime `json:"createDate"`
}

type Url struct {
	Poster     string `json:"poster"`
	Thumbnail  string `json:"thumbnail"`
	Vthumbnail string `json:"vthumbnail"`
	Hd         string `json:"hd"`
	Sd         string `json:"sd"`
}
