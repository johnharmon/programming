#!/bin/bash
#while /bin/true; do
#	if [ $1 == "hello" ]; then
#		echo "thank you"
#	elif [ $1 == "goodbye" ]; then
#		break
#	elif [ $1 == "exit" ]; then
#		echo "you exited the script"
#		exit 0
#	fi
#	echo "exit did not kill the script"
#done
#echo $1

#test=4
#echo $#
#echo $( seq 1 $# )
#echo {1.."$test"}
##echo $@
#
#
#options=$(getopt -o abcd -- "$@")
#
#while true; do
#	echo $1
#	shift
#done


#for arg in $@; do
#	#echo $arg
#	echo $1
#	shift
#done
#for name in "$*"; do 
#	echo $name
#done

#while [[ $# -gt 0 ]]; do
#	echo $#
#	echo $1
#	shift
#done

#index=0
#declare -a testarr
#
#while /bin/true; do
#	testarr[index]=${RANDOM}
#	index=$(($index+1))
#
#	if [[ $index -gt 20 ]]; then
#		break
#	fi
#done
#
#echo ${testarr[@]}







#echo $(basename $0)
#echo $(dirname $0)
#
#echo $(pwd)
VARIABLE_TEST="this value has been exported"
export $VARIABLE_TEST

