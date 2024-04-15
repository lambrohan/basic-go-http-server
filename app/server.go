package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func readFile(dir string, file string) ([]byte, error) {
	data, err := os.ReadFile(dir + "/" + file)

	if err != nil {
		return nil, err
	}

	return data, nil

}

func writeFile(dir string, file string, data string) error {
	err := os.WriteFile(filepath.Join(dir, file), []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func handleConnection(c net.Conn, directory string) {
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

	} else if strings.HasPrefix(path, "/echo") {
		suffix := path[6:]
		response = []byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\n" + "Content-Length: " + strconv.Itoa(len(suffix)) + "\r\n\r\n" + suffix)
	} else if strings.HasPrefix(path, "/user-agent") {
		userAgent := headers[2][12:]
		fmt.Println(userAgent)
		response = []byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\n" + "Content-Length: " + strconv.Itoa(len(userAgent)) + "\r\n\r\n" + userAgent)

	} else if strings.HasPrefix(path, "/files") {
		file := path[7:]
		if method == "GET" {
			data, err := readFile(directory, file)

			if err != nil {
				response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
			} else {
				response = []byte("HTTP/1.1 200 OK\r\n" + "Content-Type: application/octet-stream\r\n" + "Content-Length: " + strconv.Itoa(len(data)) + "\r\n\r\n" + string(data))
			}
		} else if method == "POST" {
			fileData := headers[len(headers)-1]
			err := writeFile(directory, file, fileData)
			if err != nil {
				response = []byte("HTTP/1.1 500 Internal Server Error\r\n\r\n")
			} else {
				response = []byte("HTTP/1.1 201 Created\r\n\r\n")
			}

		}

	} else {
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	fmt.Println(method + " " + path + " " + time.Now().Format(time.RFC3339))

	_, err = c.Write(response)
	if err != nil {
		fmt.Println("Error writing: ", err.Error())
		os.Exit(1)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	var directory string

	if len(os.Args) > 1 && os.Args[1] == "--directory" {
		directory = os.Args[2]
	}

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

		go handleConnection(connection, directory)
	}

}
