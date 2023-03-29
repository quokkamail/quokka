package main

import (
	"fmt"
	"os"

	"github.com/quokkamail/quokka/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error: %w", err))
		fmt.Fprintln(os.Stderr, "See 'quokka --help' for usage.")
		os.Exit(1)
	}
}
