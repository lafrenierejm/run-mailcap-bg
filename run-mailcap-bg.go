package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// printUsage prints a ussage message for this program to stdout then exits
// with code 1.
func printUsage() {
	fmt.Printf("usage: %s command [command_option] [...] file\n",
		os.Args[0])
	os.Exit(1)
}

// getRuntimeDir returns a path for program's runtime files that is compliant
// with the XDG Base Directory Specification.  If $XDG_RUNTIME_DIR is unset in
// the environment then os.TempDir() is used.
func getRuntimeDir(program string) string {
	var dir string

	xdgRuntime := os.Getenv("XDG_RUNTIME_DIR")

	if len(xdgRuntime) > 0 {
		dir = filepath.Join(xdgRuntime, program)
	} else {
		dir = filepath.Join(os.TempDir(), program)
	}

	return dir
}

// timestampName returns the base name of path with the current time in ISO
// 8601 format appended to it.
func timestampName(path string) string {
	// Get the current time in ISO 8601 format
	curTime := time.Now().Format("2006-01-02T15:04:05")
	// Append that to the end of path's base name
	return filepath.Base(path) + "_" + curTime
}

func main() {
	const suffix = "run-mailcap-bg"

	// Flags
	isChild := flag.Bool("child", false, "parent PID")
	flag.Parse()

	if !*isChild { // parent process
		// There must be at least 3 arguments
		numArgs := len(os.Args)
		if numArgs < 3 {
			printUsage()
		}

		// Get the directory to copy the passed file to
		dir := getRuntimeDir(suffix)
		// Attempt to create the runtime directory
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatalf("Could not create directory '%s': %v", dir, err)
		}

		// Create a timestamped copy of the file in dir
		fileOldPath := os.Args[numArgs - 1]
		fileNewPath := filepath.Join(dir, timestampName(fileOldPath))

		// Exit sucessfully
		os.Exit(0)
	} else {
		// Exit successfully
		os.Exit(0)
	}
}
