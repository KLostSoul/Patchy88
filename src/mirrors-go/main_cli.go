//go:build !windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 || (os.Args[1] != "scan" && os.Args[1] != "apply") {
		fmt.Fprintln(os.Stderr, "usage: go run . scan|apply FOLDER")
		os.Exit(2)
	}
	exeRoot := os.Getenv("MIRRORS_ASSETS_DIR")
	if exeRoot == "" {
		exeRoot = "assets"
	}
	root, err := filepath.Abs(exeRoot)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	e, err := NewEngine(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	s, err := e.Scan(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	fmt.Println("edition:", s.Edition)
	if os.Args[1] == "apply" {
		if err = e.Apply(s, func(line string) { fmt.Println(line) }); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
