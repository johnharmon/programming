#!/bin/bash

input_file="$1"
output_file="$1.converted"
converted_line=""
epoch_time=""
date_time=""
lines_converted=0
total_lines=$(wc -l $1 | cut -d ' ' -f 1)

while read line; do
	converted_line=""
	epoch_time=$(printf "${line}\n" | grep -o -E '\([0-9]{10}\.' | cut -c 2-11)
	date_time=$(date --date="@${epoch_time}")
	converted_line=$(sed -n "s/${epoch_time}\(.*\)/\(${date_time}\\1/p" <<< $line)
	[[ -n ${converted_line} ]] && line=$converted_line
	printf "$converted_line\n" >> $output_file
	((lines_converted+=1))
	#echo $lines_converted
	printf "Converted line [${lines_converted}/${total_lines}]\r"
done < ${input_file}



