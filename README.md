# GoWorkloadGenerator
Generating various workload distributions including poission and ... . **Please have a look at the other branch afterward**
## How to run
Here I will put some info on running the client and server with two configurations:

### Local
In this case you only need to run the ```myserver.go``` and ```myclient.go```. **Basically throughput and latency is not valid in this config**.

### Remote
Then in this case please take a look at the ```experiment.sh``` and ```run.sh```, it is only an experimental script for me to manage rate and number of clients. You can get the sense how I'm doing it here, however, you may have another script written for yourself to run ```client``` and ```server```remotely. 

## Commands
You may want to just run ```go run myclient.go``` fo the client and ```go run myserver.go``` for the server application.