package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
)

var total_rcv int64

func main() {

	tPTR := flag.String("type", "client", "application type")
	cmdRateIntPTR := flag.Float64("rate", 400000, "change rate of message reading")
	cmdPortPTR := flag.String("port", ":9090", "port to listen")
	clientSizePTR := flag.Int("size", 20, "number of clients")

	flag.Parse()

	t := *tPTR
	cmdRateInt := *cmdRateIntPTR
	cmdPort := *cmdPortPTR
	clientSize := *clientSizePTR

	// fmt.Println(t, cmdRateInt, cmdPort, clientSize)

	if t == "server" {
		fmt.Println("alireza")

		server(cmdPort)

	} else if t == "client" {
		t1 := time.Now()

		for i := 0; i < clientSize; i++ {
			go client(cmdRateInt, cmdPort)
		}
		// <-make(chan bool) // infinite wait.
		<-time.After(time.Second * 2)
		fmt.Println("total exchanged:", total_rcv, "\nthroughput:",
			total_rcv*1000000000/time.Now().Sub(t1).Nanoseconds(), "call/sec")

	} else if t == "client_ratelimit" {
		bucket := ratelimit.NewBucketWithQuantum(time.Second, int64(cmdRateInt), int64(cmdRateInt))
		for i := 0; i < clientSize; i++ {
			go clientRateLimite(bucket, cmdPort)
		}
		// <-make(chan bool) // infinite wait.
		<-time.After(time.Second * 3)
		fmt.Println("total exchanged", total_rcv)
	}
}

func server(cmdPort string) {
	ln, err := net.Listen("tcp", cmdPort)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go io.Copy(conn, conn)
	}
}

func client(cmdRateInt float64, cmdPort string) {

	conn, err := net.Dial("tcp", cmdPort)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}
	defer conn.Close()

	go func(conn net.Conn) {
		buf := make([]byte, 8)
		for {
			_, err := io.ReadFull(conn, buf)
			if err != nil {
				break
			}
			int_message := int64(binary.LittleEndian.Uint64(buf))
			t2 := time.Unix(0, int_message)
			fmt.Println((time.Now().UnixNano() - t2.UnixNano()) / 1000)
			atomic.AddInt64(&total_rcv, 1)
		}
		return
	}(conn)

	byte_message := make([]byte, 8)

	for {
		wait := time.Microsecond * time.Duration(nextTime(cmdRateInt)*1000000)
		if wait > 0 {
			time.Sleep(wait)
			//fmt.Println("WAIT", wait)
		}
		int_message := time.Now().UnixNano()
		binary.LittleEndian.PutUint64(byte_message, uint64(int_message))
		_, err := conn.Write(byte_message)
		if err != nil {
			log.Println("ERROR", err)
			return
		}
	}
}

func clientRateLimite(bucket *ratelimit.Bucket, cmdPort string) {

	conn, err := net.Dial("tcp", cmdPort)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}
	defer conn.Close()

	go func(conn net.Conn) {
		buf := make([]byte, 8)
		for {
			_, err := io.ReadFull(conn, buf)
			if err != nil {
				break
			}
			// int_message := int64(binary.LittleEndian.Uint64(buf))
			// t2 := time.Unix(0, int_message)
			// fmt.Println("ROUDNTRIP", time.Now().Sub(t2))
			atomic.AddInt64(&total_rcv, 1)
		}
		return
	}(conn)

	byte_message := make([]byte, 8)
	for {
		bucket.Wait(1)
		int_message := time.Now().UnixNano()
		binary.LittleEndian.PutUint64(byte_message, uint64(int_message))
		_, err := conn.Write(byte_message)
		if err != nil {
			log.Println("ERROR", err)
			return
		}
	}
}

func nextTime(rate float64) float64 {
	return -1 * math.Log(1.0-rand.Float64()) / rate
}
