#!/bin/bash

declare -A letter_array

[[ -z $* ]] && eval set -- "a"

letter_array[a]="00100 01010 01110 01010 01010"
letter_array[b]="10100 10001 10010 10001 10100"
letter_array[c]="00101 10000 10000 10000 00101"
letter_array[d]="10100 10001 10001 10001 10100"
letter_array[e]="11111 10000 11100 10000 11111"
letter_array[f]="11111 10000 11100 10000 10000"

#for item in ${letter_array[$1]}; do
    sed -n '1,${
    s/ /\n/g
    s/0/ /g
    s/1/*/g
    p
    }' <(echo ${letter_array[$1]})
#done


