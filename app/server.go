package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func handleConnection(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)
	_, err := c.Read(buf)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}

	headers := strings.Split(string(buf), "\r\n")

	reqLine := strings.Split(headers[0], " ")

	if len(reqLine) < 3 {
		fmt.Println("Invalid request")
		os.Exit(1)
	}

	path := reqLine[1]
	method := reqLine[0]

	var response []byte

	if path == "/" {
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")

	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	fmt.Println(method + " " + path + " " + time.Now().Format(time.RFC3339) + "  " + string(response))

	_, err = c.Write(response)
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
		os.Exit(1)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	fmt.Println("Listening on port 4221")

	defer l.Close()

	for {
		connection, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			os.Exit(1)
		}

		handleConnection(connection)
	}

}
