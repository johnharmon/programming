#!/bin/bash

options=$(getopt -o x --long xz, -- "$@")
eval set -- "$options"

compression='z'
extension='gz'
while /bin/true; do
	case $1 in

		-x | --xz)
			compression='J'
			extension='xz'
			;;
		--)
			break
			;;
	esac
	shift
done

datetime=$(date +%Y%m%d%H%M)
[[ -d /data/files ]] || mkdir /data/files
rm -rf /data/files/*
[[ -f /data/log ]] && rm -f /data/log
tar -cv"$compression"f /data/files/etc.tar."$extension" /etc &>> /data/log
tar -cv"$compression"f /data/files/home.tar."$extension" /home &>> /data/log
tar -cv"$compression"f /data/files/opt.tar."$extension" /opt &>> /data/log        
tar -cv"$compression"f /data/files/usr.tar."$extension" /usr &>> /data/log
#tar -cv"$compression"f /data/files/var.tar."$extension" /var &>> /data/log

mkdir /data/files/var

for directory in $(ls /var | grep -v www); do
	tar -cvf /data/files/var/$directory.tar /var/"$directory" &>> /data/log
done
tar -cv"$compression"f /data/files/var.tar."$extension" /data/files/var &>> /data/log
rm -rf /data/files/var &>> /data/log

yum list --installed > /data/files/installed_packages
date > /data/files/timestamp
tar -cvf "/data/backups/${datetime}_backup.tar" /data/files &>> /data/log
chown -R jharmon:man /data
chmod 744 -R /data
rm -rf /data/files/*

 
