package main

import (
	"fmt"
	"os"

	"engine/internal/version"
)

func main() {
	args := os.Args[1:]
	for _, a := range args {
		if a == "--" {
			continue
		}
		if a == "-v" || a == "--version" || a == "version" {
			fmt.Println(version.String())
			return
		}
	}

	fmt.Println("engine: ok")
}

