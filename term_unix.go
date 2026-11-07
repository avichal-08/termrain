//go:build !windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func stty(args ...string) (string, error) {
	c := exec.Command("stty", args...)
	c.Stdin = os.Stdin
	out, err := c.Output()
	return strings.TrimSpace(string(out)), err
}

func enableRaw() (func(), error) {
	saved, err := stty("-g")
	if err != nil {
		return nil, err
	}
	if _, err := stty("raw", "-echo"); err != nil {
		return nil, err
	}
	return func() { stty(saved) }, nil
}

func termSize() (int, int) {
	s, err := stty("size")
	if err != nil {
		return 80, 24
	}
	var rows, cols int
	if _, err := fmt.Sscan(s, &rows, &cols); err != nil {
		return 80, 24
	}
	return cols, rows
}
