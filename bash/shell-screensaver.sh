#!/bin/bash


base_string='#'

while :; do

    for color_code in {0..255}; do
        columns=$(tput cols)
        print_string=$(printf "\033[38;05;${color_code}m%${columns}s")
        echo ${print_string// /#}
        sleep 1
    done
done
     


    
