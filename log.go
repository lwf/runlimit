package main

import (
	"fmt"
	"os"
)

func fatal(reason string, format ...interface{}) {
	warning(reason, format...)
	os.Exit(1)
}

func warning(reason string, format ...interface{}) {
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", reason), format...)
}
