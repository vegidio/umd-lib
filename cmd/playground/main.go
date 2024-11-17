package main

import (
	"github.com/vegidio/umd-lib"
	"log/slog"
)

func main() {
	umdObj, _ := umd.New("https://www.reddit.com/user/atomicbrunette18/", nil, nil)
	resp, _ := umdObj.QueryMedia(10, make([]string, 0), true)
	slog.Debug("Response:", resp)
}
