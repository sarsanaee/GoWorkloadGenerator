#!/bin/bash


if [ "$#" -ne 5 ]; then
    echo "Illegal number of parameters"
    echo "./run <number_of_clients> <server_ip> <results_path> <experiments_time> <rate>"
    exit 1
fi

#result_path=$HOME/quic_results/quic
result_path=$3
rate=$5

mkdir $result_path/$rate

pkill myclient
ssh alireza@$2 "pkill my_server"
sleep 2
ssh -f alireza@$2 "export GOPATH=$HOME/work; export PATH=$PATH:/usr/local/go/bin; go run $GOPATH/src/github.com/sarsanaee/GoWorkloadGenerator/myserver.go &"


sleep 2

rm $result_path/$rate/*.log #removing current logs


for i in `seq 1 $1`;
do
	echo $i
	go run myclient.go $rate > $result_path/$rate/$i\_$1.log &
done

sleep $4



pkill myclient
ssh scc@$2 "pkill myserver"


echo "" > $result_path/$rate/latency_tcp_$rate\_$1.txt

cat $result_path/$rate/*.log > $result_path/$rate/latency_tcp_$rate\_$1.txt

