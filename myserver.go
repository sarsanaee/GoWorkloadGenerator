package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var port = ":9001"

func echo(conn net.Conn) {
	buf := make([]byte, 8) //len(message))

	for {
		_, err := io.ReadFull(conn, buf)
		if err != nil {
			return
		}
		//fmt.Println("recv")
		time.Sleep(time.Microsecond * 10)
		conn.Write(buf)
	}
}

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("ERROR", err)
			continue
		}
		go echo(conn)
	}
}
