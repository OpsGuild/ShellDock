package main

import (
	"fmt"
	"os"

	"github.com/shelldock/shelldock/internal/cli"
)

var version = "dev"

func main() {
	if err := cli.Execute(version); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
