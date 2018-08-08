#!/bin/bash

result_path=$HOME/Documents/quic_results/tcp


pkill myserver
pkill myclient

go run myserver.go &

sleep 2

for i in `seq 1 $1`;
do
	go run myclient.go > $result_path/$i.log &
done

sleep 20



pkill myserver
pkill myclient
