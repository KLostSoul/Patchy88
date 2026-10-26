//go:build !windows

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	if len(os.Args) != 3 || os.Args[1] != "apply" {
		fmt.Fprintf(os.Stderr, "usage: %s apply ORIGINAL_FOLDER\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}
	exe, _ := os.Executable()
	e, err := NewEngine(filepath.Join(filepath.Dir(exe), "assets"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	bin, cue, err := e.ApplyContext(ctx, os.Args[2], func(s string) { fmt.Println(s) })
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("완료:", bin, cue)
}
