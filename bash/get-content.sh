IFS=$'\n'
unset content
#declare -g -a content
first_element='true'
for line in $(cat $1); do
	if [[ $first_element == 'true' ]]; then 
		content=($line)
		#echo "conditional evaluated"
		first_element='false'
	else
	content+=($line)
	fi
#echo $line
done
length=${#content[@]}
#echo $length

#for index in $(seq 1 $length); do
#	echo "${content[$index]}"
	#echo -e "\r"
#done 

#export content
