package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
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
	value, err := bufio.NewReader(c).ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Println("Error:", err)
		c.Write([]byte("NOT-STORED\r\n"))
		return
	}
	key, _ := args[1], args[2]
	err = ioutil.WriteFile(key, []byte(value), 0666)
	if err != nil && err != io.EOF {
		fmt.Println("Error:", err)
		c.Write([]byte("NOT-STORED\r\n"))
		return
	}
	c.Write([]byte("STORED\r\n"))
}

func getKeyValue(c net.Conn, args []string) {
	key := args[1]
	bs, err := ioutil.ReadFile(key)
	if err != nil && err != io.EOF {
		fmt.Println("Error:", err)
	}
	c.Write([]byte("VALUE " + key + " " + strconv.Itoa(len(bs)) + " \r\n"))
	c.Write(bs)
}

func closeConnection(c net.Conn) {
	//fmt.Println("Closing connection with", c.RemoteAddr().String())
	c.Close()
}
