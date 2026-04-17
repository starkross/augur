// Command library-example demonstrates how to consume augur as a Go library.
//
//	go run ./examples/library ../bad.yaml
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/starkross/augur"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <config.yaml> [config.yaml...]", os.Args[0])
	}

	linter, err := augur.New(
		augur.WithSkipRules("OTEL-020"),
	)
	if err != nil {
		log.Fatalf("new linter: %v", err)
	}

	result, err := linter.LintFiles(context.Background(), os.Args[1:])
	if err != nil {
		log.Fatalf("lint: %v", err)
	}

	var denies, warns int
	for _, f := range result.Findings {
		marker := "WARN"
		if f.Severity == augur.SeverityDeny {
			marker = "FAIL"
			denies++
		} else {
			warns++
		}
		fmt.Printf("%s %s %s\n", marker, f.RuleID, f.Message)
	}

	fmt.Printf("\n%d failure(s), %d warning(s)\n", denies, warns)
	if denies > 0 {
		os.Exit(1)
	}
}
