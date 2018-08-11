package main

import "github.com/sarsanaee/GoWorkloadGenerator/client"

func main() {

	//runtime.GOMAXPROCS(2) // set maximum number of processes to be used by this applications

	go client.Main()

	for ture {
	}

}
