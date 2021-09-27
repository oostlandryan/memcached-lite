package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

func main() {
	// Parse command-line arguments
	portPtr := flag.Int("port", 8080, "port to listen for incoming TCP connections")
	flag.Parse()
	port := *portPtr
	fmt.Printf("Listening on port %d\n", port)
	// Listen for incoming TCP connections on specified port
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close() // close the listener when the function ends
	// For each new connection, run handleConnection in a new go-routine
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer closeConnection(c)

	for {
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error: ", err)
			}
			return
		}
		args := strings.Split(strings.TrimSpace(string(data)), " ")
		// parse command
		if args[0] == "set" {
			setKeyValue(c, args)
		} else if args[0] == "get" {
			getKeyValue(c, args)
		} else if args[0] == "EXIT" {
			break
		} else {
			c.Write([]byte("Invalid command\r\n"))
		}

	}
}

func setKeyValue(c net.Conn, args []string) {
	if len(args) < 3 {
		fmt.Println(args)
		c.Write([]byte("NOT-STORED\r\n"))
		return
	}

	key, s := args[1], args[2]
	size, err := strconv.Atoi(s)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error: ", err)
			c.Write([]byte("NOT-STORED\r\n"))
		}
		return
	}
	v := make([]byte, size+3)
	_, err = bufio.NewReader(c).Read(v)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error: ", err)
			c.Write([]byte("NOT-STORED\r\n"))
		}
		return
	}
	value := string(v[0:size])

	err = ioutil.WriteFile(key, []byte(value[0:size]), 0666)
	if err != nil {
		fmt.Println("Error:", err)
		c.Write([]byte("NOT-STORED\r\n"))
		return
	}
	c.Write([]byte("STORED\r\n"))
}

func getKeyValue(c net.Conn, args []string) {
	key := args[1]
	bs, err := ioutil.ReadFile(key)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error: ", err)
		}
		return
	}
	c.Write([]byte("VALUE " + key + " " + strconv.Itoa(len(bs)) + " \r\n"))
	c.Write([]byte(string(bs) + "\r\n"))
}

func closeConnection(c net.Conn) {
	//fmt.Println("Closing connection with", c.RemoteAddr().String())
	c.Close()
}
