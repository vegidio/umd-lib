package main

import (
	"fmt"
	"github.com/vegidio/umd-lib"
)

func main() {
	umdObj, err := umd.New("https://www.reddit.com/user/atomicbrunette18/", nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}

	resp, err := umdObj.QueryMedia(99_999, make([]string, 0), true)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Response:", resp.Media)
}
