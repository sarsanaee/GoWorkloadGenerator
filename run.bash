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
#ssh scc@$2 "pkill my_quic_server"
#sleep 2
#ssh scc@$2 "go run /home/scc/work/src/github.com/lucas-clemente/quic-go/example/echo/my_quic_server.go &"

sleep 2

for i in `seq 1 $1`;
do
	echo $i
	go run myclient.go $rate > $result_path/$rate/$i.log &
done

sleep $4



pkill myclient
ssh scc@$2 "pkill myserver"

