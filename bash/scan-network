#!/bin/bash



nmap_results=$(nmap -sn 192.168.86.0/24)

result_summary=$(echo ${nmap_results} | tail -n 1)

nmap_results=$(echo ${nmap_results} | head -n -1 | tail -n +2)

