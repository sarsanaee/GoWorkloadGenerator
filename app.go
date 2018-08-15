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

	tPtr := flag.String("type", "client", "application type")
	cmdRateIntPtr := flag.Float64("rate", 400000, "change rate of message reading")
	cmdPortPtr := flag.String("port", ":9090", "port to listen")
	clientSizePtr := flag.Int("size", 20, "number of clients")
	serverIPPtr := flag.String("ip", "10.254.254.239", "server_ip")
	expTimePtr := flag.Int("time", 5, "Experiment time")

	flag.Parse()

	t := *tPtr
	cmdRateInt := *cmdRateIntPtr
	cmdPort := *cmdPortPtr
	clientSize := *clientSizePtr
	serverIP := *serverIPPtr
	expTime := *expTimePtr

	// fmt.Println(t, cmdRateInt, cmdPort, clientSize)

	if t == "server" {
		server(cmdPort, serverIP)

	} else if t == "client" {
		t1 := time.Now()

		for i := 0; i < clientSize; i++ {
			go client(cmdRateInt, serverIP, cmdPort)
		}
		// <-make(chan bool) // infinite wait.
		<-time.After(time.Second * time.Duration(expTime))
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

func server(cmdPort string, serverIP string) {
	ln, err := net.Listen("tcp", serverIP+cmdPort)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		//go io.Copy(conn, conn)
		go echo(conn)
	}
}

func client(cmdRateInt float64, serverIP string, cmdPort string) {

	latency := ""

	conn, err := net.Dial("tcp", serverIP+cmdPort)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}

	defer conn.Close()
	defer fmt.Println(latency)

	go func(conn net.Conn) {
		buf := make([]byte, 8)
		for {
			_, err := io.ReadFull(conn, buf)
			if err != nil {
				break
			}
			// int_message := int64(binary.LittleEndian.Uint64(buf))
			// t2 := time.Unix(0, int_message)
			// fmt.Println((time.Now().UnixNano() - t2.UnixNano()) / 1000)
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
