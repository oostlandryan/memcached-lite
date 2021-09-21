package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	l, err := net.Listen("tcp4", ":9889")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

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
	//fmt.Println("Connected to ", c.RemoteAddr())
	for {
		data, err := bufio.NewReader(c).ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		args := strings.Split(strings.TrimSpace(string(data)), " ")
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
		fmt.Println("Error:", err)
		c.Write([]byte("NOT-STORED\r\n"))
		return
	}
	v := make([]byte, size+3)
	_, err = bufio.NewReader(c).Read(v)
	if err != nil && err != io.EOF {
		fmt.Println("Error:", err)
		c.Write([]byte("NOT-STORED\r\n"))
		return
	}
	value := string(v[0:size])

	err = os.WriteFile(key, []byte(value[0:size]), 0666)
	if err != nil && err != io.EOF {
		fmt.Println("Error:", err)
		c.Write([]byte("NOT-STORED\r\n"))
		return
	}
	c.Write([]byte("STORED\r\n"))
}

func getKeyValue(c net.Conn, args []string) {
	key := args[1]
	bs, err := os.ReadFile(key)
	if err != nil && err != io.EOF {
		fmt.Println("Error:", err)
	}
	c.Write([]byte("VALUE " + key + " " + strconv.Itoa(len(bs)) + " \r\n"))
	c.Write([]byte(string(bs) + "\r\n"))
}

func closeConnection(c net.Conn) {
	//fmt.Println("Closing connection with", c.RemoteAddr().String())
	c.Close()
}
