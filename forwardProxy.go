package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	// Create tcp socket which is open at port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Error creating tcp socket ", err)
	} else {
		log.Println("Listening on port 8080")
	}

	defer listener.Close()

	// Accept one connection
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting connection ", err)
		}
		// Write bytes to terminal when requests come in on it
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Read the request
		request, err := http.ReadRequest(reader)
		if err != nil {
			log.Printf("Error reading request: %v", err)
			return
		}
		log.Printf("Received request: %v", request)

		//Prepare connection to target
		host := request.URL.Host
		if !strings.Contains(host, ":") {
			host = host + ":80"
		}

		// Attempt to establish tcp connection with target
		targetConn, err := net.Dial("tcp", host)
		if err != nil {
			log.Printf("Error connecting to target %v ", err)
			return
		}

		// Forward request to target
		writeErr := request.Write(targetConn)
		if writeErr != nil {
			log.Printf("Error forwarding request to %s: %v", host, writeErr)
			targetConn.Close()
			return
		}

		response, err := http.ReadResponse(bufio.NewReader(targetConn), request)
		if err != nil {
			log.Printf("Error reading response: %v", err)
			targetConn.Close()
			return
		}
		log.Printf("Received response: %v", response)
		response.Write(conn)
		targetConn.Close()
	}
}
