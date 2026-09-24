//go:build !windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if (len(os.Args) != 3 && len(os.Args) != 4) || (os.Args[1] != "scan" && os.Args[1] != "apply") {
		fmt.Fprintln(os.Stderr, "usage: go run . scan|apply FOLDER [Japanese|English]")
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
	preferred := ""
	if len(os.Args) == 4 {
		switch strings.ToLower(os.Args[3]) {
		case "japanese", "japan", "jp":
			preferred = "Japanese"
		case "english", "eng", "en":
			preferred = "English"
		default:
			fmt.Fprintln(os.Stderr, "edition must be Japanese or English")
			os.Exit(2)
		}
	}
	s, err := e.ScanWithEdition(os.Args[2], preferred)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if len(s.Options) > 1 && s.Edition == "" {
		fmt.Println("available:", strings.Join(s.Options, ", "))
		fmt.Println("To choose, pass Japanese or English after FOLDER")
		if os.Args[1] == "apply" {
			os.Exit(2)
		}
		return
	}
	fmt.Println("edition:", s.Edition)
	if os.Args[1] == "apply" {
		if err = e.Apply(s, func(line string) { fmt.Println(line) }); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
