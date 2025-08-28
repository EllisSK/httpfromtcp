package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Failed to read messages.txt with error:", err)
	}

	lines := getLinesChannel(file)

	for line := range lines {
		fmt.Println("read:", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				ch <- str
				str = ""
			}
			str += string(data)
		}

		if len(str) != 0 {
			ch <- str
		}
	}()

	return ch
}
