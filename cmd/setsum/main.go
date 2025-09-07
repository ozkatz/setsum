package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ozkatz/setsum"
)

func usage() {
	fmt.Println("Commands:")
	fmt.Println("insert <value>")
	fmt.Println("remove <value>")
	fmt.Println("merge <hex-digest>")
	fmt.Println("subtract <hex-digest>")
	fmt.Println("digest")
}

func main() {
	// read lines from stdin
	scanner := bufio.NewScanner(os.Stdin)
	ss := setsum.Default()
	for scanner.Scan() {
		line := scanner.Text()

		if line == "digest" {
			fmt.Println(ss.HexDigest())
			continue
		}
		// line begins with "insert ":
		if strings.HasPrefix(line, "insert ") {
			// remove "insert " from the line
			line = strings.TrimPrefix(line, "insert ")
			ss.Insert([]byte(line))
			continue
		}
		if strings.HasPrefix(line, "remove ") {
			line = strings.TrimPrefix(line, "remove ")
			ss.Remove([]byte(line))
			continue
		}
		if strings.HasPrefix(line, "merge ") {
			line = strings.TrimPrefix(line, "merge ")
			ss = ss.Merge(setsum.FromHexDigest(line))
			continue
		}
		if strings.HasPrefix(line, "subtract ") {
			line = strings.TrimPrefix(line, "subtract ")
			ss = ss.Subtract(setsum.FromHexDigest(line))
			continue
		}
		// if line is one of q, quit, exit:
		if strings.HasPrefix(line, "q") || strings.HasPrefix(line, "quit") || strings.HasPrefix(line, "exit") {
			break
		}
		usage()
	}
}
