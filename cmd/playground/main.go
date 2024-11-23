package main

import (
	"encoding/json"
	"fmt"
	"github.com/vegidio/umd-lib"
)

func main() {
	umdObj, err := umd.New("https://www.reddit.com/user/atomicbrunette18/", nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}

	resp, err := umdObj.QueryMedia(99999, make([]string, 0), true)
	if err != nil {
		fmt.Println("Error:", err)
	}

	j, _ := json.MarshalIndent(resp.Media, "", "  ")
	fmt.Println("Response:", string(j))
}
