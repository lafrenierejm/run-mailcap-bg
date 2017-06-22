package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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

// Author: markc (https://stackoverflow.com/a/21067803)
//
// copyFileContents copies the contents src to dst. dst will be created if it
// does not already exist. If dst does exist it will be replaced by the
// contents of src.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
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

		// Copy contents of fileOrigPath to fileNewPath
		err = copyFileContents(fileOldPath, fileNewPath)
		if err != nil {
			log.Fatalf("Error copying '%s' to '%s': %v",
				fileOldPath, fileNewPath, err)
		}
		// Set the new path as argument for the child process
		os.Args[numArgs-1] = fileNewPath

		// Get the absolute path of this executable
		cmdPath, err := filepath.Abs(os.Args[0])
		if err != nil {
			log.Fatalf("Error resolving path of executable: %v", err)
		}

		// Create the child process
		os.Args[0] = "-child" // cmdPath is used in place of Args[0]
		cmd := exec.Command(cmdPath, os.Args...)
		err = cmd.Start()
		if err != nil {
			log.Fatalf("Command finished with error: %v", err)
		}

		// Exit sucessfully
		os.Exit(0)
	} else {
		// Args[0] is the name of this executable.
		// Args[1] is "-child"
		cmd := exec.Command(os.Args[2], os.Args[3:]...)
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Command finished with error: %v", err)
		}

		// The last argument was the file which can now be removed.
		file := os.Args[len(os.Args) - 1]
		err = os.Remove(file)
		if err != nil {
			log.Fatalf("Unable to remove file '%s': %v", file, err)
		}

		// Exit successfully
		os.Exit(0)
	}
}
