package main

import (
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/vegidio/umd-lib"
)

func main() {
	array := []int{1, 2, 3, 4, 5, 6}
	filter := []int{1, 2}

	newArray := funk.FilterInt(array, func(value int) bool {
		return !funk.ContainsInt(filter, value)
	})

	fmt.Printf("Array: %v\n", newArray)

	umdObj, _ := umd.New("https://www.reddit.com/user/atomicbrunette18/", make(map[string]interface{}), nil)
	resp := umdObj.QueryMedia(100, make([]string, 0))
	fmt.Printf("Size: %d", len(resp.Media))
}
