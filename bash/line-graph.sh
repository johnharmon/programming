#!/bin/bash
#set -xv
x=0
y=0
c=1
length=10
current_x=0
current_y=0
line=""
num_lines=$length
result=''
neg_coefficient=1

while [[ $# -gt 0 ]]; do
	case $1 in

		-x ) 
			shift
			x=$1
			shift
			;;
		-y ) 
			shift
			y=$1
			shift
			;;
		-c | --coefficient )
			shift
			c=$1
            if [[ $c -lt 0 ]]; then c=${c%-}; neg_coefficient=0; fi
			shift
			;;
		-l | --length )
			shift
			length=$1
			shift
			;;
		* )
			printf "\"$1\" is not an accepted argument\n"
			printf "Accepted arguments are:\n"
			printf '%s\n' "-x"
			printf '%s\n' "-y"
			printf '%s\n' "-c or --coefficient"
			printf '%s\n' "-l or --length"
			exit 1
			;;
	esac
done

#echo "x = $x"
current_x=$x
current_y=$y

printf "x=$x\n"
printf "y=$y\n"
printf "coe=$c\n"	
printf "length=$length\n"
printf "current_x=$current_x\n"

get_num_lines (){
	localx=$(bc <<< "scale=2; $x+$length-1")
	local local_lines=$(bc <<< "scale=2; $c*$localx+$y")
        #echo "localx =  $localx"
		((local_lines++))
	echo "${local_lines#-}"
}
#       echo $x $current_x

evaluation (){
	local result=$(bc <<< "scale=2; $c*$current_x+$y")
    result=${result#-}
	echo $result
}			

make_charlines (){
	local line=""
	for i in $(seq 1 $length); do
		line="$line+"
	done
	echo $line
}

make_output_file (){
	printf "" > /tmp/graph
#	sed -i '1 s/.*/graph/g' /tmp/graph
}

populate_output_file (){
	local line=$1
	echo "num_lines = $num_lines"
	for new_line in $(seq 1 $num_lines); do
		#sed  -i "$new_line a $line" /tmp/graph
		printf "%s\n" "$line" >> /tmp/graph
	done
}

make_output_file
num_lines=$(get_num_lines)
line=$(make_charlines)
echo $line
populate_output_file $line
result=$(evaluation)

while [[ $current_x -le $(bc <<< "scale=2; $x+$length") ]]; do
	result=$(evaluation)
	echo "y=$result, x=$current_x"
	((result++))
	sed -i "${result}s/\(^+\{$current_x\}\)+\(.*\)/\1*\2/g" /tmp/graph
	#sed -i "${result}s/\(^+\{$current_x\}\)\(.*\)/\1*\2/g" /tmp/graph
	((current_x++))
done

if [[ $neg_coefficient -eq 1 ]]
    then

    ending_line=$num_lines
    while [[ $ending_line -gt 0 ]]; do
        sed -n "${ending_line}p" /tmp/graph >> /tmp/reversed_graph
        printf "%s\n" "$ending_line"
        ((ending_line--))
    done < /tmp/graph
     
    rm -f /tmp/graph
    mv /tmp/reversed_graph /tmp/graph

fi







		 






