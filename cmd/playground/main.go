package main

import (
	"fmt"
	"github.com/vegidio/umd-lib"
)

func main() {
	umdObj, _ := umd.New("https://www.reddit.com/user/atomicbrunette18/", nil, nil)
	resp, _ := umdObj.QueryMedia(10, make([]string, 0), true)
	fmt.Printf("Response: %v\n", resp.Media)
}
