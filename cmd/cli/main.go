package main

import (
	"bibgo"
	"fmt"
)

func main() {
    // TODO add tree walking
    // TODO add elapsed time
    // TODO add logger
	counter := bibgo.ParseFile("data/input/acm/acm.bib")
	fmt.Printf("Parsed %d entries\n", counter)
}
