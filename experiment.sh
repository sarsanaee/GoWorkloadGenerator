#!/bin/bash

for i in `seq 1 5 63`;
do
	./run.sh $i 10.254.254.239 ~/quic_results/tcp 30 5000
done
