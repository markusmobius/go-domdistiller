// +build ignore

package main

import (
	"fmt"

	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	result, err := distiller.ApplyForFile("example/sample.html", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.HTML)
}
