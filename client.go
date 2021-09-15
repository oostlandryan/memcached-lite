package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for j := 0; j < 1; j++ {
		for i := 0; i < 500; i++ {
			wg.Add(1)
			fileNum := i + j*1000
			go testMemcache(fileNum, &wg)
		}
		wg.Wait()
	}
}

func testMemcache(i int, wg *sync.WaitGroup) {
	con, err := net.Dial("tcp", "localhost:9889")
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	defer wg.Done()
	time.Sleep(time.Second)
	serverReader := bufio.NewReader(con)
	filename := "firstTest" + strconv.Itoa(i)
	con.Write([]byte("set " + filename + ".testfile 0 \r\n"))
	time.Sleep(500 * time.Millisecond)
	con.Write([]byte("this is my " + strconv.Itoa(i) + "th test string\r\n"))
	time.Sleep(500 * time.Millisecond)

	// Waiting for the server response
	serverResponse, err := serverReader.ReadString('\n')
	//fmt.Println("Received server response")
	switch err {
	case nil:
		fmt.Println(strings.TrimSpace(serverResponse))
	case io.EOF:
		fmt.Println("server closed the connection")
		return
	default:
		fmt.Printf("server error: %v\n", err)
		return
	}

	filename = "secondTest" + strconv.Itoa(i)
	con.Write([]byte("set " + filename + ".testfile 0 \r\n"))
	time.Sleep(500 * time.Millisecond)
	con.Write([]byte("this is my " + strconv.Itoa(i) + "th test string\r\n"))
	time.Sleep(500 * time.Millisecond)
	con.Write([]byte("EXIT\r\n"))

	// Waiting for the server response
	serverResponse, err = serverReader.ReadString('\n')
	//fmt.Println("Received server response")
	switch err {
	case nil:
		fmt.Println(strings.TrimSpace(serverResponse))
	case io.EOF:
		fmt.Println("server closed the connection")
		return
	default:
		fmt.Printf("server error: %v\n", err)
		return
	}
}
