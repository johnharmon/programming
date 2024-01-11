#!/bin/bash


output_sting=""




for num in $(seq 0 255 ); do
#	color_code=$(bc <<< ${num}%7+1)
	printf "Counting: \033[38;5;${num}m${num} \033[00m \r"
	sleep 1
done
