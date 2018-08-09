#!/bin/bash

for i in `seq 1 5 90`;
do
	./run.bash $i 10.254.254.1 ~/quic_results/tcp 100 10000
done
