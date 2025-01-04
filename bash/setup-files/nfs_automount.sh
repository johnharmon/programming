#!/bin/bash

[[ -d /nfs ]] || mkdir /nfs
[[ -d /nfs/home ]] || mkdir /nfs/home
[[ -d /nfs/data ]] || mkdir /nfs/data
[[ -d /nfs/etc ]] || mkdir /nfs/etc

ping 192.168.137.2 -c 1 && {
	mount -t nfs 192.168.137.2:/home /nfs/home
	mount -t nfs 192.168.137.2:/data /nfs/data
	mount -t nfs 192.168.137.2:/etc /nfs/etc
}

exit 0

