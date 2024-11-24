package main

import (
	"encoding/json"
	"fmt"
	"github.com/vegidio/umd-lib"
)

func main() {
	umdObj, err := umd.New("https://coomer.su/onlyfans/user/belledelphine/post/1101526004", nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
	}

	resp, err := umdObj.QueryMedia(99999, make([]string, 0), true)
	if err != nil {
		fmt.Println("Error:", err)
	}

	j, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println("Response:", string(j))
	fmt.Println("Number of files:", len(resp.Media))
}
