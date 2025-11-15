package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// telnet commands
const (
	IAC  = 255
	DONT = 254
	DO   = 253
	WONT = 252
	WILL = 251
)

func main() {
	var timeout int
	flag.IntVar(&timeout, "timeout", 10, "Connection timeout, standard - 10s")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Printf("Usage: %s [--timeout=10] <host> <port>\n", os.Args[0])
		fmt.Println("Examples:")
		fmt.Printf("  %s --timeout=5 google.com 80\n", os.Args[0])
		fmt.Printf("  %s localhost 23\n", os.Args[0])
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	fmt.Printf("Connecting to %s with timeout %d \n", address, timeout)

	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Fatalf("Error connecting %s: %v", address, err)
	}
	defer conn.Close()

	fmt.Printf("Successfully connected %s\n", address)

	done := make(chan bool, 2)

	go readFromSocketWithTelnet(conn, done)
	go writeToSocket(conn, done)

	idleTimeout := time.After(time.Duration(timeout) * time.Second)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalCh:
		fmt.Println("\nStopping")
	case <-done:
		fmt.Println("Connection has been cut")
	case <-idleTimeout:
		fmt.Println("\nIdle timeout - connection has been closed")
	}
}

func readFromSocketWithTelnet(conn net.Conn, done chan<- bool) {
	reader := bufio.NewReader(conn)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			log.Printf("Error reading from sockets: %v", err)
			done <- true
			return
		}

		if n > 0 {
			processedData := processTelnetCommands(buffer[:n])
			if len(processedData) > 0 {
				fmt.Print(string(processedData))
			}
		}
	}
}

func processTelnetCommands(data []byte) []byte {
	result := make([]byte, 0, len(data))
	i := 0

	for i < len(data) {
		if data[i] == IAC && i+2 < len(data) {
			i += 3
		} else {
			result = append(result, data[i])
			i++
		}
	}

	return result
}

func writeToSocket(conn net.Conn, done chan<- bool) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input := scanner.Text() + "\r\n"

		_, err := conn.Write([]byte(input))
		if err != nil {
			log.Printf("Error writing into socket: %v", err)
			done <- true
			return
		}
	}

	if err := scanner.Err(); err != nil {
		if err == io.EOF {
			fmt.Println("\nFinishing...")
		} else {
			log.Printf("Error reading enter: %v", err)
		}
	}
	done <- true
}
