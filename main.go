package main

import (
	"fmt"
	"my-go-db/server"
	"os"
	"strconv"
	"strings"
	"flag"
	"bufio"
)

const listSep = ","
const keyValueSep = ":"

var prompt string


func printPromt() {
	fmt.Print(prompt)
}

func parseToStringSlice(input string) ([]string, bool) {
	if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		input = strings.TrimLeft(input, "[")
		input = strings.TrimRight(input, "]")
		stringSlice := strings.Split(input, listSep)
		return stringSlice, true
	}
	return nil, false
}

func parseToIntSlice(input string) ([]int, bool) {
	stringSlice, ok := parseToStringSlice(input)
	if !ok {
		return nil, false
	}

	res := []int{}
	for _, s := range stringSlice {
		if intValue, err := strconv.Atoi(s); err != nil {
			return nil, false
		} else {
			res = append(res, intValue)
		}
	}
	return res, true
}

func parseToStringMap(input string) (map[string]string, bool) {
	if strings.HasPrefix(input, "{") && strings.HasSuffix(input, "}") {
		input = strings.TrimLeft(input, "{")
		input = strings.TrimRight(input, "}")
		pairs := strings.Split(input, listSep)
		stringMap := make(map[string]string)
		for _, p := range pairs {
			kv := strings.Split(p, keyValueSep)
			if len(kv) == 2 {
				stringMap[strings.Trim(kv[0], " ")] = strings.Trim(kv[1], " ")
			}
		}
		if len(stringMap) > 1 {
			return stringMap, true
		}
	}
	return nil, false
}

func parseToIntMap(input string) (map[string]int, bool) {
	stringMap, ok := parseToStringMap(input)
	if !ok {
		return nil, false
	}

	intMap := make(map[string]int)
	for k, v := range stringMap {
		if intValue, err := strconv.Atoi(v); err != nil {
			return nil, false
		} else {
			intMap[k] = intValue
		}
	}
	return intMap, true
}

func CMD_KEYS(c *server.Client) {
	keys := c.GetKeys()
	fmt.Printf("%v\n", keys)
}

func CMD_SET(c *server.Client, key, input string, ttl int) {
	// SET INT
	if intValue, err := strconv.Atoi(input); err == nil {
		if err := c.SetInt(key, intValue, ttl); err != nil {
			fmt.Println("Error: ", err.Error())
		}
		return
	}

	// SET INT SLICE
	if intSliceValue, ok := parseToIntSlice(input); ok {
		if err := c.SetIntSlice(key, intSliceValue, ttl); err != nil {
			fmt.Println("Error: ", err.Error())
		}
		fmt.Println("SET INT")
		return
	}

	// SET STRING SLICE
	if stringSliceValue, ok := parseToStringSlice(input); ok {
		if err := c.SetStringSlice(key, stringSliceValue, ttl); err != nil {
			fmt.Println("Error: ", err.Error())
		}
		return
	}

	// SET INT MAP
	if intMap, ok := parseToIntMap(input); ok {
		if err := c.SetIntMap(key, intMap, ttl); err != nil {
			fmt.Printf("Error: ", err.Error())
		}
		return
	}

	// SET STRING MAP
	if stringMap, ok := parseToStringMap(input); ok {
		if err := c.SetStringMap(key, stringMap, ttl); err != nil {
			fmt.Printf("Error: ", err.Error())
		}
		return
	}

	// SET STRING
	if err := c.SetString(key, input, ttl); err != nil {
		fmt.Println("Error: ", err.Error())
	}
}

func CMD_GET(c *server.Client, key string) {
	if value, err := c.GetValue(key); err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Printf("%v\n", value)
	}
}

func CMD_REMOVE(c *server.Client, key string) {
	if err := c.Remove(key); err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("Done")
	}
}


func printUsage() {
	fmt.Println("my-go-db is a tool to run server or client to server")
}


var host string
var port string


func init() {
	flag.StringVar(&host, "host", "localhost", "Server's host")
	flag.StringVar(&port, "port", "8080", "Server's port")
	flag.Parse()
}


func startServer() {
	fmt.Println("Starting server on port", port)
	s := server.New(":" + port)
	s.Start()
	s.WaitStop()
}


func startClient() {
	fmt.Printf("Connecting to server http://%s:%s\n", host, port)

	prompt = fmt.Sprintf("%s:%s$ >>> ", host, port)


	client := server.NewClient(host, port)

	scanner := bufio.NewScanner(os.Stdin)
	printPromt()
	for scanner.Scan() {
		input := strings.Split(scanner.Text(), " ")
		switch len(input) {
		case 1:
			cmd := input[0]
			switch strings.ToUpper(cmd) {
			case "KEYS":
				CMD_KEYS(client)
			}
		case 2:
			cmd, key := input[0], input[1]
			switch strings.ToUpper(cmd) {
			case "GET":
				CMD_GET(client, key)
			case "REMOVE":
				CMD_REMOVE(client, key)
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
				CMD_SET(client, key, value, ttl)
			}
		}
		printPromt()
	}
	if scanner.Err() != nil {
		fmt.Printf("Error: %v", scanner.Err())
	}
}


func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "server":
		startServer()
	case "client":
		startClient()
	default:
		printUsage()
	}
}
