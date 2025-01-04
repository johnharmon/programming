#!/bin/bash


while :; do
    for directory in $(find /home/jharmon -type d); do
        clear
        printf "Contents of ${directory}:\n"
        ls ${directory}
        sleep 2
    done
done


