package main

import (
	"fmt"
	"os"
	"log"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Failed to read messages.txt with error: %v", err)
	}

	for {
		data := make([]byte, 8)
		part, err := file.Read(data)
		if err != nil {
			break
		}
		fmt.Printf("read: %s\n", string(data[:part]))
	}
}
