package main

import (
	"bufio"
	"fmt"
	"log"
	"my-go-db/server"
	"os"
	"strconv"
	"strings"
)

func main() {

	// todo start server if need
	fmt.Println("Starting server...")
	s := server.New(":8080")
	s.Start()

	log.Println("Connecting to server...")
	cl := server.NewClient("localhost", "8080")

	fmt.Println("Put something...")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := strings.Split(scanner.Text(), " ")
		switch len(input) {
		case 1:
			cmd := (input[0])
			switch strings.ToUpper(cmd) {
			case "KEYS":
				continue
			}
		case 2:
			cmd, key := input[0], input[1]
			switch strings.ToUpper(cmd) {
			case "GET":
				value, _ := cl.GetValue(key)
				fmt.Println(value)
			case "REMOVE":
				continue
			}
		case 3, 4:
			cmd, key, value := input[0], input[1], input[2]
			ttl := 0
			if len(input) == 4 {
				ttl2, err := strconv.Atoi(input[3])
				if err != nil {
					fmt.Println("Bad ttl value. Must be integer")
					continue
				}
				ttl = ttl2
			}
			switch strings.ToUpper(cmd) {
			case "SET":
				value, _ := cl.SetValue(key, value, int64(ttl))
				fmt.Println(value)
			}
		}
	}
	if scanner.Err() != nil {
		fmt.Printf("Error: %v", scanner.Err())
	}
}
