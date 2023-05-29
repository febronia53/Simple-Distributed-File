package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.4:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	tcpConn := conn.(*net.TCPConn)
	defer tcpConn.Close()

	fmt.Println("Connected to server")

	// Send a request to the server to get the aggregated character count map
	request := "GET_AGGREGATED_CHAR_COUNT"
	_, err = tcpConn.Write([]byte(request))
	if err != nil {
		fmt.Println(err)
		return
	}
	tcpConn.CloseWrite()

	// Receive the aggregated character count map from the server
	var responseBuf bytes.Buffer
	_, err = io.Copy(&responseBuf, conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	dec := gob.NewDecoder(&responseBuf)
	aggregatedCharCounts := make(map[string]int)
	err = dec.Decode(&aggregatedCharCounts)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write the aggregated character count map to the file
	file, err := os.Create("received_file.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for char, count := range aggregatedCharCounts {
		_, err = file.WriteString(fmt.Sprintf("%s: %d\n", char, count))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Aggregated data received and stored in received_file.txt")
}
