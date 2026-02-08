package main

import (
	"fmt"
	"os"
)

func printHelp() {
	fmt.Println("gh-problemas: A GitHub CLI extension for managing problemas")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  gh problemas [command]")
	fmt.Println("")
	fmt.Println("Available Commands:")
	fmt.Println("  version     Show version information")
	fmt.Println("  help        Show this help message")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	command := os.Args[1]

	switch command {
	case "version", "-v", "--version":
		fmt.Println("gh-problemas version 0.1.0")
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Run 'gh problemas help' for usage information")
		os.Exit(1)
	}
}
