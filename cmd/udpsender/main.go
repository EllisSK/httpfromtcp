package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	UDPaddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		fmt.Println("Error in UDP resolver!")
		log.Fatal(err)
	}

	UDPconn, err := net.DialUDP("udp", nil, UDPaddr)
	defer UDPconn.Close()

	inReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		line, err := inReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error in reader!")
			log.Fatal(err)
		}

		_, err = UDPconn.Write([]byte(line))
		if err != nil {
			fmt.Println("Error in writer!")
			log.Fatal(err)
		}
	}
}
