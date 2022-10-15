package main

import (
	"fmt"

	"github.com/Dogancan94/toolkit"
)

func main() {
	var tools toolkit.Tools
	s := tools.CreateRandomString(10)

	fmt.Println("Random string:", s)
}
