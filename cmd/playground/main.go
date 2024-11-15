package main

import (
	"fmt"
	"github.com/vegidio/umd-lib"
	"github.com/vegidio/umd-lib/event"
)

func main() {
	umdObj, _ := umd.New("https://www.redgifs.com/watch/liquidlostgoldfish", nil, func(event event.Event) {
		fmt.Printf("Event: %v\n", event)
	})

	resp, _ := umdObj.QueryMedia(100, make([]string, 0))
	fmt.Printf("Response: %v\n", resp)
}
