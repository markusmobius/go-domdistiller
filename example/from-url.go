// +build ignore

package main

import (
	"fmt"
	"time"

	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	url := "https://arstechnica.com/gadgets/2020/10/iphone-12-and-12-pro-double-review-playing-apples-greatest-hits/"

	// Start distiller
	result, err := distiller.ApplyForURL(url, time.Minute, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.HTML)
}
