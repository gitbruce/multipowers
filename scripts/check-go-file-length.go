package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	max := 500
	violations := 0
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		f, e := os.Open(path)
		if e != nil {
			return nil
		}
		defer f.Close()
		s := bufio.NewScanner(f)
		lines := 0
		for s.Scan() {
			lines++
		}
		if lines > max {
			fmt.Printf("WARN %s has %d lines (max %d)\n", path, lines, max)
			violations++
		}
		return nil
	})
	if violations > 0 {
		os.Exit(1)
	}
}
