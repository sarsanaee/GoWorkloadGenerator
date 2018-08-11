package main

import (
	"os"
	"strconv"

	"github.com/sarsanaee/GoWorkloadGenerator/client"
)

func main() {

	//runtime.GOMAXPROCS(2) // set maximum number of processes to be used by this applications
	args := os.Args[1:]
	rate_int, _ := strconv.Atoi(args[0])
	client_size, _ := strconv.Atoi(args[1])
	port := args[2]

	i := 0
	for i <= client_size {
		go client.Main(rate_int, port)
		i = i + 1
	}

	for true {

	}

}
