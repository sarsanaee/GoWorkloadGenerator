package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	mrand "math/rand"
	"net"
	"os"
	"time"
)

var port = "0.0.0.0:9001"

func nextTime(rate float64) float64 {
	return -1 * math.Log(1.0-mrand.Float64()) / rate
}

func main() {

        args := os.Args[1:]
        rate_int, err := strconv.Atoi(args[0])
	var rate float64 = float64(rate_int)

	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	var my_random_number float64 = nextTime(rate) * 1000000
	var my_random_int int = int(my_random_number)
	var int_message int64 = time.Now().UnixNano()
	byte_message := make([]byte, 8)

	go func(conn net.Conn) {
		buf := make([]byte, 8)

		for true {
			_, err = io.ReadFull(conn, buf)
			now := time.Now().UnixNano()

			if err != nil {
				return
			}

			last := int64(binary.LittleEndian.Uint64(buf))

			fmt.Println((now - last) / 1000)
		}
		return

	}(conn)

	for true {
		time.Sleep(time.Microsecond * time.Duration(my_random_int))
		int_message = time.Now().UnixNano()
		binary.LittleEndian.PutUint64(byte_message, uint64(int_message))
		conn.Write(byte_message)
	}
}
