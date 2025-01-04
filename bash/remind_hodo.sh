#!/bin/bash


reminded_hodo=0
today=$(date +"%A")

if [[ ${today} == "Wednesday" ]] && [[ ${reminded_hodo} -eq 0 ]]; then
    reminded_hodo=1
    popup_message "Give Kyle and example of where he isn't logical"
elif [[ ${today} != "Wednesday" ]]; then
    reminded_hodo=0
fi



