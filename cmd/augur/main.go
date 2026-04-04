package main

import "os"

var version = "dev"

func main() {
	if err := newRootCmd(version).Execute(); err != nil {
		os.Exit(1)
	}
}
