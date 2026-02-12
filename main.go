package main

import (
	"os"

	"github.com/hpg/gh-problemas/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
