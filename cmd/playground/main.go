package main

import (
	"github.com/vegidio/kmd-lib/internal/extractors/reddit"
)

func main() {
	reddit := reddit.Reddit{}
	reddit.QueryMedia("https://www.reddit.com/user/atomicbrunette18", 10, nil)
}
