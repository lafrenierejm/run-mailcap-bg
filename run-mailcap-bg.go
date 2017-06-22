package main

import (
	"flag"
	"fmt"
	"os"
)

// printUsage prints a ussage message for this program to stdout then exits
// with code 1.
func printUsage() {
	fmt.Printf("usage: %s command [command_option] [...] file\n",
		os.Args[0])
	os.Exit(1)
}

func main() {
	// Flags
	isChild := flag.Bool("child", false, "parent PID")
	flag.Parse()

	if !*isChild { // parent process
		// There must be at least 3 arguments
		numArgs := len(os.Args)
		if numArgs < 3 {
			printUsage()
		}

		// Exit sucessfully
		os.Exit(0)
	} else {
		// Exit successfully
		os.Exit(0)
	}
}
