#!/bin/bash

for i in `seq 30 1 50`;
do
	./run.sh $i 10.254.254.239 ~/quic_results/tcp 30 1000
done
