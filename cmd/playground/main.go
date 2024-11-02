package main

import "github.com/vegidio/umd-lib/internal/extractors/reddit"

func main() {
	reddit.GetUserSubmissions("atomicbrunette18", "", 100)
}
