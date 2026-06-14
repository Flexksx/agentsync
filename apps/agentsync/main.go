package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		_, _ = fmt.Fprintln(os.Stdout, "agentsync — sync skills and instructions across AI agent vendors")
		os.Exit(0)
	}
	_, _ = fmt.Fprintln(os.Stdout, "agentsync — sync skills and instructions across AI agent vendors")
}
