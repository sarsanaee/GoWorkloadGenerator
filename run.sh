#!/bin/bash

if [ "$#" -ne 5 ]; then
    echo "Illegal number of parameters"
    echo "./run <number_of_clients> <server_ip> <results_path> <experiments_time> <rate>"
    exit 1
fi

size=$1
server_ip=$2
result_path=$3
time=$4
rate=$5

mkdir $result_path/$rate

pkill myclient
ssh alireza@server_ip "pkill app"
sleep 2
ssh -f alireza@server_ip "export GOPATH=$HOME/work; export PATH=$PATH:/usr/local/go/bin; go run $GOPATH/src/github.com/sarsanaee/GoWorkloadGenerator/myserver.go &"

sleep 2

#rm $result_path/$rate/*.log #removing current logs

go run app.go -size $size -ip $server_ip -type client -rate $rate -time $time > $result_path/$rate/$size.log 

ssh scc@$server_ip "pkill app"

sleep(2)


