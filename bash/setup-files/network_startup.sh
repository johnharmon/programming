#!/bin/bash

#ip addr | grep "192.168.0"

#res=$?

#if [[ $res -ne 0 ]]; then
#	nmcli con down eth0
#	nmcli con up private
#fi
#nmcli con up eth0
dhclient -r
dhclient -v --timeout 10
#ping -c 1 8.8.8.8 

#res=$?

#if [[ $res -ne 0 ]]; then
#	nmcli con down eth0
#	nmcli con up private
#fi
nmcli con up private

