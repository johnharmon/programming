#!/bin/bash

stripLeadingZero (){
	if [[ "${nanoseconds:0:1}" == '0' ]]; then
		nanoseconds=${nanoseconds:1}
	fi
}

nanoseconds=$(date +%N)
stripLeadingZero


eSeconds=$(date +%s)
seed=$(($eSeconds * $nanoseconds / ${RANDOM}))
seed2=${eSeconds}${nanoseconds}
echo $seed
echo $seed2

