package coomer

import (
	"github.com/vegidio/umd-lib/internal/utils"
)

type Response struct {
	Post   *Post  `json:"post"`
	Images []File `json:"previews"`
	Videos []File `json:"attachments"`
}

type Post struct {
	Id        string         `json:"id"`
	Service   string         `json:"service"`
	User      string         `json:"user"`
	Published utils.NotzTime `json:"published"`
}

type File struct {
	Server string `json:"server"`
	Path   string `json:"path"`
}
