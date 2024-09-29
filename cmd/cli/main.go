package main

import (
	"bibgo"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// TODO add logger
	start := time.Now()
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
	elapsed := float64(time.Since(start).Microseconds()) / 1000
	average := elapsed / float64(counter)
	fmt.Printf("Took     %7.3f ms to parse %d entries\n", elapsed, counter)
	fmt.Printf("Averaged %7.3f ms per entry\n", average)
}
