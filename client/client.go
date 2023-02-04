package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func connect(url string, id int, okchan chan int, ch chan int) {

	conn, err := net.Dial("tcp", url)
	if err != nil {
		// handle error
		fmt.Println("ACH PROBLEM {} {}", err, conn)
		os.Exit(1)
	}
	defer conn.Close()
	okchan <- id
	time.Sleep(2 * time.Second)
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		// handle error
		fmt.Println("ACH PROBLEM {} {}", err, conn)
		os.Exit(1)
	}
	fmt.Printf("Got %d\n", result[0])

}

func main() {
	var url = os.Args[1]
	fmt.Printf("Will connect to %s\n", url)
	time.Sleep(1 * time.Second)
	ch := make(chan int)
	okchan := make(chan int)

	for i := 0; i < 30000; i++ {
		go connect(url, i, okchan, ch)
		var count = <-okchan
		if count%100 == 0 {
			fmt.Printf("Connection %d ok\n", count)
		}
	}

	fmt.Printf("connections are open, press enter to quit...")
	reader := bufio.NewReader(os.Stdin)
	_, _, err := reader.ReadRune()
	if err != nil {
		// handle error
		fmt.Println("Cant scan PROBLEM {}", err)
		os.Exit(1)
	}
	fmt.Printf("Quitting\n")

}
