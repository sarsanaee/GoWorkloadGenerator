#!/bin/bash

if [ "$#" -ne 3 ]; then
    echo "Illegal number of parameters"
    echo "./run <statrt> <increament> <stop> <rate>"
    exit 1
fi


for i in `seq $1 $2 $3`;
do
	./run.sh $i 10.254.254.239 ~/quic_results/tcp 30 1000
done

python3 throughput_latency.py 1000 300000 $1 $3 $2 ~/quic_results/tcp/1000
