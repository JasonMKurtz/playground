package main

import (
	"fmt"
	"os"
	"strings"
)

func tail(filename string, numlines int) {
	file, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	size := file.Size()
	f, _ := os.Open(filename)

	bytesize := 1024
	bytes := make([]byte, bytesize)
	f.Read(bytes)

	for i := size; i >= 0; i -= int64(bytesize) {
		f.Seek(i, 2)
		content := string(bytes)
		lines := strings.Split(content, "\n")

		start := len(lines) - numlines - 1
		if start < 0 {
			start = 0
		}

		fmt.Printf("%s", strings.Join(lines[start:], "\n"))
	}

	f.Close()
}

func main() {
	tail("file", 1)
}
