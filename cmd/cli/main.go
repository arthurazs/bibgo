package main

import (
	"bibgo"
	"fmt"
)

func main() {
	counter := bibgo.ParseFile("data/input/acm/acm.bib")
	fmt.Printf("Parsed %d entries\n", counter)
}
