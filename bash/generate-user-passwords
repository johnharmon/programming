#!/bin/bash

#for arg in $@; do
#	echo $arg
#	shift
#done

genChar (){
	code=$(bc <<< "${RANDOM}%94+33")
	char=$(echo $code | awk '{ printf("%c",$0); }')
	echo "${char}"
}

genPass (){
	length=16
	for char in $(seq 1 $length); do
	password+=$(genChar)	
	done
	echo "${password}"
}



userpass=''

for arg in $@; do
        userpass=$(genPass)
	echo "${arg}"
        echo "${userpass}"
done

