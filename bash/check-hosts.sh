#!/bin/bash

while read host; do
    host_name=$(echo "$host" | awk '{print $1}')
    ip_address=$(echo "$host" | awk '{print $3}')
    printf "pinging $host_name @ $ip_address:\n"
    (ping $host_name;)
    printf "\033[2J\033[0;0H"
done <<< $(cat /var/named/mydomain.com.zone | grep 192)



#for host in $(cat /var/named/mydomain.com.zone | grep 192); do
#    host_name=$(echo "$host" | awk '{print $1}')
#    ip_address=$(echo "$host" | awk '{print $3}')
#    printf "pinging $host_name @ $ip_address:\n"
#    (ping $host_name;)
#done

