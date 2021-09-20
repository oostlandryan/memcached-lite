package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type memcacheConnection struct {
	con net.Conn
}

func main() {
	var wg sync.WaitGroup
	for j := 0; j < 1; j++ {
		for i := 0; i < 5; i++ {
			wg.Add(1)
			fileNum := i + j*1000
			go testSetGet(fileNum, &wg)
		}
		wg.Wait()
	}
}

// Creates a new memcacheConnection to the given address
func newMemcacheConnection(address string) (memcacheConnection, error) {
	con, err := net.Dial("tcp4", address)
	if err != nil {
		return memcacheConnection{}, err
	}
	return memcacheConnection{con: con}, nil
}

// Closes memcacheConnection's TCP connection
func (m memcacheConnection) Close() {
	m.con.Close()
}

// Method of memcacheConnection to call the set command
func (m memcacheConnection) set(key string, value string) (string, error) {
	// Write to server
	m.con.Write([]byte("set " + key + " " + strconv.Itoa(len(value)) + " \r\n"))
	time.Sleep(250 * time.Millisecond)
	m.con.Write([]byte(value + " \r\n"))
	// Get server response
	serverReader := bufio.NewReader(m.con)
	serverResponse, err := serverReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return serverResponse, nil

}

// Method of memcacheConnection to call the get command
func (m memcacheConnection) get(key string) (string, error) {
	// Write to server
	m.con.Write([]byte("get " + key + " \r\n"))
	// Get server response
	serverReader := bufio.NewReader(m.con)
	serverResponse, err := serverReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	responseSlice := strings.Split(serverResponse, " ")
	var size int
	size, err = strconv.Atoi(responseSlice[2])
	if err != nil {
		return "", err
	}
	size = size - 3
	// Read data block
	result := make([]byte, size)
	_, err = serverReader.Read(result)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// Sets a value and checks if it gets the same value
func testSetGet(i int, wg *sync.WaitGroup) {
	memCon, err := newMemcacheConnection("localhost:9889")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	defer wg.Done()
	defer memCon.Close()

	var response string
	key := "key" + strconv.Itoa(i) + ".ryoost"
	value := "This is the " + strconv.Itoa(i) + "th value"
	response, err = memCon.set(key, value)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println(response)

	var gotValue string
	gotValue, err = memCon.get(key)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Printf("Set: %s Get: %s\n", "\""+value+"\"", "\""+gotValue+"\"")

}
