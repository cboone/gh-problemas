package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("gh-problemas: A GitHub CLI extension for managing problemas")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  gh problemas [command]")
		fmt.Println("")
		fmt.Println("Available Commands:")
		fmt.Println("  version     Show version information")
		fmt.Println("  help        Show this help message")
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "version", "-v", "--version":
		fmt.Println("gh-problemas version 0.1.0")
	case "help", "-h", "--help":
		fmt.Println("gh-problemas: A GitHub CLI extension for managing problemas")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  gh problemas [command]")
		fmt.Println("")
		fmt.Println("Available Commands:")
		fmt.Println("  version     Show version information")
		fmt.Println("  help        Show this help message")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'gh problemas help' for usage information")
		os.Exit(1)
	}
}
