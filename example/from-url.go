// +build ignore

package main

import (
	"fmt"
	"time"

	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	url := "https://www.vice.com/en/article/k7qpqe/how-coronavirus-is-impacting-the-arab-world"

	// Start distiller
	result, err := distiller.ApplyForURL(url, time.Minute, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.HTML)
}
