package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var connCount = 0

func handleConn(conn net.Conn, ch chan int, connId int) {
	if connId%100 == 0 {
		fmt.Printf("Connection count: %d - remote is %s\n", connId, conn.RemoteAddr().String())
		if connId > 65536 {
			fmt.Printf("Quand est-ce que je visite le bunker???\n")
		}
	}
	<-ch
	fmt.Fprintf(conn, "%d", connId)
}

func vazy(ln net.Listener, ch chan int) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("ACH PROBLEM {}", err)

		}
		connCount++
		var connId = connCount
		go handleConn(conn, ch, connId)
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
	go vazy(ln, ch)

	// on attend un pressage de bouton
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		// handle error
		fmt.Println("Cant scan PROBLEM {}", err)
		os.Exit(1)
	}
	fmt.Printf("Sending end to all connections, %d\n", char)
	for i := 0; i < connCount; i++ {
		ch <- i
	}

}
