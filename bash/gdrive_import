#!/bin/bash
trap "unset IFS; echo ""; exit 100" SIGINT


IFS=$'\n'
format='text/plain'
cd /home/jharmon/gdrive_local
for file in $(ls -l /home/jharmon/gdrive_local); do
    printf "%s\n" "$file"
    read -p "Would you like to upload this document? [y/n]: "  key
    if [[ ${key} == "y" ]]; then 
        file_name=$(echo $file | awk '{print $9}')
        echo $file_name
        gdrive import --mime ${format} ${file_name} 
    fi
done


unset IFS

