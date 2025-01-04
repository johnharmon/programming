#!/bin/bash

trap "unset ITS; echo ""; exit 100" SIGINT


IFS=$'\n'
format='text/plain'
for file in $(gdrive list); do
    printf "%s\n" "$file"
    read -p "Would you like to download this document? [y/n]: "  key
    if [[ ${key} == "y" ]]; then 
        file_id=$(echo $file | awk '{print $1}')
        echo $file_id
        gdrive export --mime ${format} ${file_id} 
    fi
done
unset IFS

