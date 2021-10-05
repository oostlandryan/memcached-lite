package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var serverAddress string

type memcacheConnection struct {
	con net.Conn
}

func main() {
	// Parse command-line arguments
	flag.StringVar(&serverAddress, "server", "localhost:8080", "address server is listening on")
	flag.Parse()
	fmt.Println("------------------------------------")
	fmt.Printf("Client will connect to %s\n", serverAddress)
	fmt.Println("------------------------------------")
	// Run Tests
	testMassConcurrency(500)
	//testKeyNames()
	//testLargeValue()
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
	// Read data block
	result := make([]byte, size)
	_, err = serverReader.Read(result)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// Sets a value and checks if it gets the same value
func testSetGet(i int, c chan bool) {
	memCon, err := newMemcacheConnection(serverAddress)
	if err != nil {
		fmt.Println("Error: ", err)
		c <- false
		return
	}

	defer memCon.Close()

	var response string
	key := "key" + strconv.Itoa(i) + ".ryoost"
	value := "This is the " + strconv.Itoa(i) + "th value at %s" + time.Now().String()
	response, err = memCon.set(key, value)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	if response == "NOT-STORED\r\n" {
		c <- false
		return
	}

	var gotValue string
	gotValue, err = memCon.get(key)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	c <- value == gotValue

}

// Calls testSetGet connections number of times and prints out the number of success and failures
func testMassConcurrency(connections int) {
	c := make(chan bool, connections)
	for i := 0; i < connections; i++ {
		go testSetGet(i, c)
	}
	correctCount := 0
	incorrectCount := 0
	for i := 0; i < connections; i++ {
		result := <-c
		if result {
			correctCount++
		} else {
			incorrectCount++
		}
	}
	fmt.Println("---------------------------------------------------------------------------------------")
	fmt.Printf("Test Mass Concurrency: %d correct get-set operations, %d incorrect get-set operations\n", correctCount, incorrectCount)
	fmt.Println("---------------------------------------------------------------------------------------")
}

// Tests if large (>250 chars) keys and keys with special characters can work
func testKeyNames() {
	memCon, err := newMemcacheConnection(serverAddress)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer memCon.Close()
	var longKey string
	key := "A"
	for i := 0; i < 247; i++ {
		key = key + "A"
	}
	var response string
	response, err = memCon.set(key+".ryoost", "The key is "+key)
	if err != nil {
		fmt.Println("Error: ", err)
		longKey = "failed"
	}
	if response == "NOT-STORED\r\n" {
		longKey = "failed"
	} else {
		longKey = "passed"
	}

	var tooLongKey string
	key = "A"
	for i := 0; i < 248; i++ {
		key = key + "A"
	}
	response, err = memCon.set(key+".ryoost", "The key is "+key)
	if err != nil {
		fmt.Println("Error: ", err)
		tooLongKey = "failed"
	}
	if response == "NOT-STORED\r\n" {
		tooLongKey = "failed"
	} else {
		tooLongKey = "passed"
	}

	var newlineKey string
	key = "my\nstrange\rkey"
	response, err = memCon.set(key+".ryoost", "The key is "+key)
	if err != nil {
		fmt.Println("Error: ", err)
		newlineKey = "failed"
	}
	if response == "NOT-STORED\r\n" {
		newlineKey = "failed"
	} else {
		newlineKey = "passed"
	}

	var spaceKey string
	key = "my strange key"
	response, err = memCon.set(key+".ryoost", "The key is "+key)
	if err != nil {
		fmt.Println("Error: ", err)
		spaceKey = "failed"
	}
	if response == "NOT-STORED\r\n" {
		spaceKey = "failed"
	} else {
		spaceKey = "passed"
	}

	var specialKey string
	key = "my&strange|key"
	response, err = memCon.set(key+".ryoost", "The key is "+key)
	if err != nil {
		fmt.Println("Error: ", err)
		specialKey = "failed"
	}
	if response == "NOT-STORED\r\n" {
		specialKey = "failed"
	} else {
		specialKey = "passed"
	}

	fmt.Printf("254 byte key: %s\n", longKey)
	fmt.Printf("255 byte key: %s\n", tooLongKey)
	fmt.Printf("Newline Key: %s\n", newlineKey)
	fmt.Printf("Key with spaces: %s\n", spaceKey)
	fmt.Printf("Key with special characters: %s\n", specialKey)
}

// Tests if large values (~4KB) can be stored and retrieved
func testLargeValue() {
	memCon, err := newMemcacheConnection(serverAddress)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer memCon.Close()
	value := "A"
	for i := 0; i < 4000; i++ {
		value = value + "A"
	}
	key := "largeFile.ryoost"
	var response string
	response, err = memCon.set(key, value)
	if err != nil {
		fmt.Println("Error: ", err)
		fmt.Println("Large Value: failed")
		return
	}
	if response == "NOT-STORED\r\n" {
		fmt.Println("Large Value: failed")
		return
	}

	var gotValue string
	gotValue, err = memCon.get(key)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	if gotValue != value {
		fmt.Println("Large Value: failed")
	} else {
		fmt.Println("Large Value: passed")
	}
}
