package main

import (
	"os"

	"github.com/fin-assistant/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
