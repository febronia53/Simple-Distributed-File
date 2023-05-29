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
	// Read the file to send
	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Count character occurrences
	charCount := make(map[rune]int)
	buffer := make([]byte, 1024)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		if bytesRead == 0 {
			break
		}
		for _, char := range string(buffer[:bytesRead]) {
			charCount[char]++
		}
	}

	conn, err := net.Dial("tcp", "192.168.1.4:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	tcpConn := conn.(*net.TCPConn)
	defer tcpConn.Close()

	fmt.Println("Connected to server")

	// Send the character count map
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(charCount)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = tcpConn.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	tcpConn.CloseWrite()

	fmt.Println("Character count sent successfully")
	conn.Close()
}
