package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var connCount = 0

func handleConn(conn net.Conn, ch chan int, connId int, wg *sync.WaitGroup) {
	wg.Add(1)
	if connId%1000 == 0 {
		fmt.Printf("Connection count: %d - remote is %s\n", connId, conn.RemoteAddr().String())
		if connId > 65536 {
			fmt.Printf("Quand est-ce que je visite le bunker???\n")
		}
	}
	for {
		select {
		case <-ch:
			// finished, we send a connId
			fmt.Fprintf(conn, "%d\n", connId)
			fmt.Printf("Finishing connection %d -- remote was %s\n", connId, conn.RemoteAddr().String())
			wg.Done()
		case <-time.After(5 * time.Second):
			// send a ping regularly
			fmt.Fprintln(conn, ".")
		}
	}

}

func vazy(ln net.Listener, ch chan int, wg *sync.WaitGroup) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("ACH PROBLEM {}", err)

		}
		connCount++
		var connId = connCount
		go handleConn(conn, ch, connId, wg)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
		fmt.Println("ACH PROBLEM {}", err)
		os.Exit(1)
	}
	fmt.Printf("Server listening on %s\n", ln.Addr().String())

	ch := make(chan int)
	var wg sync.WaitGroup
	go vazy(ln, ch, &wg)

	// on attend un pressage de bouton
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		// handle error
		fmt.Println("Cant scan PROBLEM {}", err)
		os.Exit(1)
	}
	fmt.Printf("Sending end to all %d connections, %d\n", connCount, char)
	for i := 0; i < connCount; i++ {
		ch <- i
	}
	fmt.Println("Waiting on all connections to complete...")
	wg.Wait()
	fmt.Println("Done, seeya")

}
