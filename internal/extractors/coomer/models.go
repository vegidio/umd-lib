package coomer

import (
	"github.com/vegidio/umd-lib/internal/utils"
)

type Post struct {
	Service     string         `json:"service"`
	User        string         `json:"user"`
	Published   utils.NotzTime `json:"published"`
	File        File           `json:"file"`
	Attachments []File         `json:"attachments"`
}

type Response struct {
	Post *Post `json:"post"`
}

type File struct {
	Path string `json:"path"`
}
