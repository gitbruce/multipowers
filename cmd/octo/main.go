package main

import (
	"os"

	"github.com/gitbruce/claude-octopus/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
