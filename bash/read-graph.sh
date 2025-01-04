#!/bin/bash


while read line; do

	formatted_line=$(sed 's/*/\\033\[01;31m*\\033\[00;32m/g' <<< $line)
	printf "$formatted_line\n"

done < /tmp/graph


