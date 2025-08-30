package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		fmt.Println("Error in listener!")
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error in conn!")
			log.Fatal(err)
		}
		fmt.Println("A connection has been accepted.")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error in request!")
			log.Fatal(err)
		}

		fmt.Println("Request line:")
		fmt.Println("- Method:", req.RequestLine.Method)
		fmt.Println("- Target:", req.RequestLine.RequestTarget)
		fmt.Println("- Version:", req.RequestLine.HttpVersion)

		fmt.Println("The connection has been closed")
	}
}
