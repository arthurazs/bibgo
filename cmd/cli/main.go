package main

import (
	"bibgo"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// TODO add elapsed time
	// TODO add logger
	counter := uint64(0)

	path := filepath.Join("data", "input")
	dirs, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	fmt.Println("Opening folders...")
	for _, dir := range dirs {
		subFilepath := filepath.Join(path, dir.Name())
		fmt.Printf("Opening %s...\n", subFilepath)
		subdirs, err := os.ReadDir(subFilepath)
		if err != nil {
			panic(err)
		}
		for _, dir = range subdirs {
			counter += bibgo.ParseFile(filepath.Join(subFilepath, dir.Name()))
		}
	}
	fmt.Printf("Took     ? ms to parse %d entries\n", counter)
}
