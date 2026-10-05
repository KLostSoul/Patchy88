//go:build !windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 || os.Args[1] != "apply" {
		fmt.Fprintf(os.Stderr, "usage: %s apply ORIGINAL.zip\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}
	exe, _ := os.Executable()
	e, err := NewEngine(filepath.Join(filepath.Dir(exe), "assets"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	out, err := e.Apply(os.Args[2], func(s string) { fmt.Println(s) })
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("완료:", out)
}
