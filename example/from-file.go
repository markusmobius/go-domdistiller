// +build ignore

package main

import (
	"fmt"

	"github.com/go-shiori/dom"
	distiller "github.com/markusmobius/go-domdistiller"
)

func main() {
	result, err := distiller.ApplyForFile("example/sample.html", nil)
	if err != nil {
		panic(err)
	}

	rawHTML := dom.OuterHTML(result.Node)
	fmt.Println(rawHTML)
}
