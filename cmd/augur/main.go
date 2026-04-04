package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	if err := newRootCmd(version).Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
