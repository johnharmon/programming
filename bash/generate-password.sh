#!/bin/bash


password=''

for number in {1..170}; do


password+=$(bc <<< "${RANDOM}%93+33" | awk '{ printf("%c",$0); }')

done

echo $password
