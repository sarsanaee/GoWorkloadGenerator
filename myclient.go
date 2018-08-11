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

	go client.Main(rate_int)

	for true {

	}

}
