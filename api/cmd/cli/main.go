package main

import (
	"os"
)

func main() {
	root, cleanup := newRoot()
	defer cleanup()

	if err := run(root); err != nil {
		os.Exit(1)
	}
}
