#!/bin/bash
arg1=''
arg2=''
arg3=''

while getopts ":ab:c" opt; do
	case ${opt} in 
	a )
	#	echo $opt
		arg1='option -a was invoked'
		;;
	b )
	#	echo $opt
		if [ "${OPTARG:0:1}" != '-' ]; then 
			arg2=$OPTARG
		else
			echo "Option b must be given a non-parameter argument"

		fi
		;;
	c )
	#	echo $opt
		if [ $OPTARG ]; then
			arg3=$OPTARG
		else
			arg3='no value was passed to -c'
	
		fi
		;;
	: )
		echo "Invalid option: $opt requires an argument"
		;;
	\? )
		echo $opt
		echo "invalid option specified $opt" ;;

	esac
done

#echo "arguments: -a: $arg1; -b: $arg2; -c: $arg3"


if [ ! -z "$arg1" ]; then
	echo $arg1
fi

if [ ! -z "$arg2" ]; then
	echo $arg2
fi

if [ ! -z "$arg3" ]; then
	echo $arg3
fi

#echo $arg1
#echo $arg2
#echo $arg3
