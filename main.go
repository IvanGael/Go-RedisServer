package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type RedisServer struct {
	data map[string]string
}

func NewRedisServer() *RedisServer {
	return &RedisServer{
		data: make(map[string]string),
	}
}

func (s *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		parts := strings.Fields(command)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "SET":
			if len(parts) != 3 {
				fmt.Fprintf(conn, "Error: SET command requires key and value\n")
				continue
			}
			key := parts[1]
			value := parts[2]
			s.data[key] = value
			fmt.Fprintf(conn, "OK\n")
		case "GET":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "Error: GET command requires key\n")
				continue
			}
			key := parts[1]
			value, ok := s.data[key]
			if !ok {
				fmt.Fprintf(conn, "(nil)\n")
				continue
			}
			fmt.Fprintf(conn, "%s\n", value)
		case "DEL":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "Error: DEL command requires key\n")
				continue
			}
			key := parts[1]
			delete(s.data, key)
			fmt.Fprintf(conn, "OK\n")
		default:
			fmt.Fprintf(conn, "Error: unknown command '%s'\n", parts[0])
		}
	}
}

func main() {
	server := NewRedisServer()
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Redis server started on port 6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go server.handleConnection(conn)
	}
}

//Enable Telnet on windows
//open powershell as administrator
// Enable-WindowsOptionalFeature -FeatureName TelnetClient -Online

//Example usage
// telnet localhost 6379
