package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func connect(url string, id int, okchan chan int, wg *sync.WaitGroup) {

	wg.Add(1)
	conn, err := net.Dial("tcp", url)
	if err != nil {
		// handle error
		fmt.Println("ACH PROBLEM {} {}", err, conn)
		os.Exit(1)
	}
	defer conn.Close()
	okchan <- id

	reader := bufio.NewReader(conn)
	var pingCount = 0
	for {
		result, err := reader.ReadString('\n')
		if err != nil {
			// handle error
			fmt.Println("ACH PROBLEM {} {}", err, conn)
			wg.Done()
			return
		}
		switch result {
		case ".\n":
			//fmt.Println(".")
			pingCount++
		default:
			fmt.Printf("Conn %d finished -- %d pings\n", id, pingCount)
			wg.Done()
			return
		}
	}

}

func main() {
	var url = os.Args[1]
	fmt.Printf("Will connect to %s\n", url)
	time.Sleep(1 * time.Second)
	okchan := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < 25000; i++ {
		go connect(url, i, okchan, &wg)
		var count = <-okchan
		if count%1000 == 0 {
			fmt.Printf("Connection %d ok\n", count)
		}
	}

	fmt.Println("Connections launched, waiting for them to complete")
	wg.Wait()
	fmt.Println("Done, quitting")

}
