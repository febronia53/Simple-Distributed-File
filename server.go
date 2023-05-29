package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	_ "github.com/go-sql-driver/mysql"
)

func handleConnection(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	// Read the request from the client
	request := make([]byte, 1024)
	n, err := conn.Read(request)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Client disconnected")
		} else {
			fmt.Println(err)
		}
		return
	}

	request = request[:n] // Trim the request to the actual length

	if string(request) == "GET_AGGREGATED_CHAR_COUNT" {
		// Query the database to get the character counts
		rows, err := db.Query("SELECT `char`, SUM(count) FROM char_counts GROUP BY `char`")

		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()

		// Create a map to store the character counts
		charCounts := make(map[string]int)

		// Iterate over the rows and store the character counts in the map
		for rows.Next() {
			var char string
			var count int
			err := rows.Scan(&char, &count)
			if err != nil {
				fmt.Println(err)
				return
			}
			charCounts[char] = count
		}

		// Encode the map and send it to the client
		var responseBuf bytes.Buffer
		enc := gob.NewEncoder(&responseBuf)
		err = enc.Encode(charCounts)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.Write(responseBuf.Bytes())
	} else {
		// Handle the character count map sent by the slave
		dec := gob.NewDecoder(bytes.NewReader(request))

		charCount := make(map[rune]int)
		err = dec.Decode(&charCount)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Store the map in the database
		stmt, err := db.Prepare("INSERT INTO char_counts (`char`, count) VALUES (?, ?) ON DUPLICATE KEY UPDATE count = count + ?")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer stmt.Close()

		for char, count := range charCount {
			if char != ' ' && char != '\n' && char != '\r' { // Skip spaces and newline characters
				_, err = stmt.Exec(string(char), count, count)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func main() {
	// Connect to the database
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/count")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// Create table
	_, err = db.Exec("CREATE TABLE  IF NOT EXISTS char_counts (id INT NOT NULL AUTO_INCREMENT, `char` VARCHAR(1) NOT NULL, count INT NOT NULL, PRIMARY KEY (id));")
	if err != nil {
		panic(err.Error())
	}

	// Modify column
	_, err = db.Exec("ALTER TABLE char_counts MODIFY `char` VARCHAR(1) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin;")
	if err != nil {
		panic(err.Error())
	}

	// Add unique index
	_, err = db.Exec("ALTER TABLE char_counts ADD UNIQUE (`char`);")
	if err != nil {
		panic(err.Error())
	}
	// Listen for incoming connections on port 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	// Wait for incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn, db)
	}
}
