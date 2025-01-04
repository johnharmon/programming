#!/bin/bash


username=''
fullname=''
password=''

genChar (){
	code=$(bc <<< "${RANDOM}%94+33")
	char=$(echo $code | awk '{ printf("%c",$0); }')
	echo "${char}"
}

genPass (){
	password=''
	length=16
	for char in $(seq 1 $length); do
	password+=$(genChar)	
	done
	echo "${password}"
}


if [[ ${#} -eq 0 ]]; then
	echo "Usage: add-local-user.sh [username]. [comments]..." 1>&2
	exit 1
else

	uid=$(id -u)
	if [[ $uid -ne 0 ]]; then
		echo "Not being execuded with sudo privleges. Exiting..." 1>&2
		exit 2
	fi

	username=$1
	shift
	fullname="$*"

	userpass=$(genPass)
	useradd $username -p "${userpass}" -c "${fullname}" 1>/dev/null
	res=$?

	if [[ $res -ne 0 ]]; then
		echo "Failed to add user, the useradd command exited with the following error code: $res" 1>&2
		exit 3
	fi

	

	echo "username: $username"
	echo "full name: $fullname"
	echo "password: $userpass"
	echo "Hostname: $(hostname)"

	passwd -e $username 1>/dev/null
fi



#read -p "Please enter a username: " username
#read -p "What is the full name of the user for this account? " fullname
#read -p "Please enter a password: " password
#
#userpass=''
#
#for arg in $@; do
#        userpass=$(genPass)
#	echo "${arg}"
#        echo "${userpass}"
#done


#uid=$( /home/jharmon/c_programs/compiled_programs/getuid)
