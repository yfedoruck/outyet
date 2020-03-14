package main

import (
	"os"
	"path/filepath"
)

func parentPath() string {
	dir, _ := os.Getwd()
	return dir + filepath.FromSlash("/src/github.com/yfedoruck/outyet/")
}
