package main

import (
	"os"
	"log"
	"strings"
	"path/filepath"
)

func visit(path string, f os.FileInfo, err error) error {
	if ! strings.HasSuffix(path, ".html") { return nil }

	log.Printf("Visited: %s\n", path)		
	return nil
} 

func main() {
	ptt_dir := os.Args[1]

	log.Printf("arg: %s", ptt_dir)
	err := filepath.Walk(ptt_dir, visit)
	log.Printf("filepath.Walk() returned %v\n", err)
}
