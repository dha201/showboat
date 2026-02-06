package main

import (
	"fmt"
	"os"

	"github.com/simonw/showcase/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		if len(os.Args) < 4 {
			fmt.Fprintln(os.Stderr, "usage: showcase init <file> <title>")
			os.Exit(1)
		}
		if err := cmd.Init(os.Args[2], os.Args[3]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "build":
		fmt.Fprintln(os.Stderr, "build: not yet implemented")
		os.Exit(1)
	case "verify":
		fmt.Fprintln(os.Stderr, "verify: not yet implemented")
		os.Exit(1)
	case "extract":
		fmt.Fprintln(os.Stderr, "extract: not yet implemented")
		os.Exit(1)
	case "--help", "-h", "help":
		printUsage()
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`showcase - Create executable demo documents that show and prove an agent's work.

Showcase helps agents build markdown documents that mix commentary, executable
code blocks, and captured output. These documents serve as both readable
documentation and reproducible proof of work. A verifier can re-execute all
code blocks and confirm the outputs still match.

Usage:
  showcase init <file> <title>             Create a new demo document
  showcase build <file> commentary [text]  Append commentary (text or stdin)
  showcase build <file> run <lang> [code]  Run code and capture output
  showcase build <file> image [script]     Run script, capture image output
  showcase verify <file> [--output <new>]  Re-run and diff all code blocks
  showcase extract <file>                  Emit build commands to recreate file

Global Options:
  --workdir <dir>   Set working directory for code execution (default: current)
  --help, -h        Show this help message

Stdin:
  The build subcommands accept input from stdin when the text/code argument is
  omitted. For example:
    echo "Hello world" | showcase build demo.md commentary
    cat script.sh | showcase build demo.md run bash

Example:
  # Create a demo
  showcase init demo.md "Setting Up a Python Project"

  # Add commentary
  showcase build demo.md commentary "First, let's create a virtual environment."

  # Run a command and capture output
  showcase build demo.md run bash "python3 -m venv .venv && echo 'Done'"

  # Run Python and capture output
  showcase build demo.md run python "print('Hello from Python')"

  # Capture a screenshot
  showcase build demo.md image "python screenshot.py http://localhost:8000"

  # Verify the demo still works
  showcase verify demo.md

  # See what commands built the demo
  showcase extract demo.md

Resulting markdown format:

  # Setting Up a Python Project

  *2026-02-06T15:30:00Z*

  First, let's create a virtual environment.

  ` + "```" + `bash
  python3 -m venv .venv && echo 'Done'
  ` + "```" + `

  ` + "```" + `output
  Done
  ` + "```" + `

  ` + "```" + `python
  print('Hello from Python')
  ` + "```" + `

  ` + "```" + `output
  Hello from Python
  ` + "```" + `
`)
}
