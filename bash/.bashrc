# .bashrc
#set -o vi
# Source global definitions
#set -xv
if [ -f /etc/bashrc ]; then
	. /etc/bashrc
fi
#cd ~
# User specific environment
PATH="$HOME/.local/bin:$HOME/bin:$PATH"
PATH="/home/jharmon/bash-scripts:$PATH"
export PATH
export PYTHONPATH="PYTHONPATH:/home/jharmon/.local/lib/python3.9/site-packages"
profile='/home/jharmon/.bashrc'
export profile
MYVIMRC='/home/jharmon/.vimrc'
export MYVIMRC
ANSIBLE_PATH='/home/jharmon/programming/ansible/work'
# Uncomment the following line if you don't like systemctl's auto-paging feature:
# export SYSTEMD_PAGER=

# User specific aliases and functions

alias get-content='source /home/jharmon/bashscripts/get-content.sh'
alias em='set +o vi; set -o emacs'
alias vi='set +o emacs; set -o vi'
alias vim='vim -u /home/jharmon/.vimrc'
alias sm='ps -ef | grep -v grep | grep -E'
alias rm='rm -i'
alias env-path='env | grep PATH | cut -d'=' -f 2 | tr ":" "\n"'
alias centos2='ssh 192.168.137.2'
alias ld='find . -maxdepth 1 -type d'
CENTOS2='192.168.137.2'
export CENTOS2
export HISTTIMEFORMAT="%h %d %H:%M:%S "
export HISTFILESIZE=1000000
export HISTSIZE=1000000
shopt -s histappend
#PROMPT_COMMAND="$PROMPT_COMMAND; history -a"
export HISTCONTROL=ignorespace
[[ ${SUDO_USER} == 'jharmon' ]] && touch /home/jharmon/.root_history && HISTFILE='/home/jharmon/.root_history'
#function get-content {

#IFS=$'\n'

#declare -g -a result

#for line in $(cat $1); do
#	result+=($line)
#echo $line
#done

#length=${#result[@]}
#echo $length

#f#or index in $(seq 1 $length); do
#	echo "${result[$index]}"
	#echo -e "\r"
#done 

#export result 
#}

killjobs (){

	for job in $(jobs -l | tr -s ' ' | cut -d' ' -f2); do
		kill -9 $job
	done


}
cdu (){

	up_string=''
	up_levels=1
	next_path=""
	up_levels=$(echo ${1} | cut -d'/' -f1)
#	printf "before string creation\n"	
	#grep -E '[0-9]{1,2}/[^[:space:]].*' <<< $1 1>/dev/null && next_path=$(echo $1 | cut -d'/' -f2- | grep -v '^-')
	grep -E '[0-9]{1,2}/.*' <<< $1 1>/dev/null && next_path=$(echo $1 | cut -d'/' -f2- | grep -v '^-')
#	printf "after string creation\n"	
	for level in $(seq 1 $up_levels); do 
		up_string+='../'
	done
	result_string=${up_string}${next_path}
	cd ${result_string}
}
#PS1='\[\033[38;5;199m\]j\[\033[00;32m\]h\[\033[00;33m\]a\[\033[00;34m\]r\[\033[00;35m\]m\[\033[00;36m\]o\[\033[00;32m\]n\[\033[00;31m\]@$HOSTNAME:$PWD>\[\033[00;32m\] '
#change-color() {
#
#	options=$(getopt -o "" --long color:,text: -- "$@" )
#
#	echo $options
#	eval set -- $options
#
#	colorset='31'
#	output=''
#
#
#	while /bin/true; do
#		case $1 in
#			--color)
#				shift
#				echo $1
#				if [[ $1 == 'red' ]] || [[ $1 == 'Red' ]] || [[ $1 == 'RED' ]];
#				then
#					colorset='31'
#				elif [[ $1 == 'blue' ]] || [[ $1 == 'Blue' ]] || [[ $1 == 'BLUE' ]];
#				then
#					colorset='34'
#				fi
#				;;
#			--text)
#				shift
#				output=$1
#				;;
#			--)
#				shift
#				break
#				;;
#		esac
#		shift
#	done
#
#
#	PS1='$USER@$HN:$PWD> '
#	PS2='continue->'
#
#
#}


show-colors (){ 
	options=$(getopt -o ls:t:o: --longoptions list,string:,trail:,stdout:,short,format -- "$@")
	eval set -- "$options"

	standard_output='Color #'
	list="false"
	trailing_invoked='false'
	output_string=""
	trailing_sequence=""
	while /bin/true; do
		case $1 in

		-l | --list)
			list="true"
			;;
		-s | --string)
			shift
			output_string="$1"
			;;
		-t | --trail)
			shift
			trailing_sequence="$1"
			trailing_invoked='true'
			;;
		-o | --stdout)
			shift
			standard_output="$1"
			;;
		--short) 
			short="true"
			 ;;
		--format)
			format="true"
			;;
		--)
			shift
			break
			;;
		esac
		shift
	done
	if [[ -n $short ]]; then
		printf "\033[01;30m#30(Black) \033[01;31m#31(Red) \033[01;32m#32(Green)  \033[01;33m#33(Yellow) \033[01;34m#34(Blue) \033[01;35m#35(Magenta) \033[01;36m#36(Cyan) \033[01;37m#37(White)\033[00m\n"
		unset short
		return

	elif [[ -n $format ]]; then
		printf "Formatting for standard bash colorization of text and background is as follows:\n"
		 echo 'All color codes must be preceeded by \033[<attribute_code>;<background_code>;<foreground_code>m'
		 printf "(Note: Multiple attribute codes can be used, simply separate them via semicolons like the rest of the codes)\n"
		 printf "Foreground Color codes are:\n"
		printf "\033[01;30m#30(Black) \033[01;31m#31(Red) \033[01;32m#32(Green)  \033[01;33m#33(Yellow) \033[01;34m#34(Blue) \033[01;35m#35(Magenta) \033[01;36m#36(Cyan) \033[01;37m#37(White)\033[00m\n"
		printf "Background Color codes are:\n"
		printf "\033[01;40;37m#40(Black) \033[01;41;37m#41(Red) \033[01;42;37m#42(Green)  \033[01;43;37m#43(Yellow) \033[01;44;37m#44(Blue) \033[01;45;37m#45(Magenta) \033[01;46;37m#46(Cyan) \033[01;47;31m#47(White)\033[00m\n"
		printf "Attribute modifiers are as follows:\n"
		printf "\033[00m#00 for Normal\n\033[01m#01 for Bold\n\033[04m#04 for Underscore\n\033[05;41m#05 for Blinking background\n\033[00m\033[07m#07 for Reversed\n\033[08m#08 for Concealed \033[00m\n"
		printf "The last one is #08 for concealed incase your terminal actually reads the codes properly\n"
		unset format
		return
	
	else
		if [[ $list == "true" ]]; then
			trailing_sequence='\n'
		elif [[ -z ${trailing_sequence} ]]; then
			trailing_sequence='\t'
		fi


			for num in $(seq 1 255); do
				printf "\033[38;5;${num}m${standard_output}${num} ${output_string}${trailing_sequence}"
			done | fold -s -w 240

			printf "\n"
	fi
}

test-function (){
	
	for arg; do echo $arg; done

}

set +xv

isvalidip(){
    grep -E '^[0-255]\.[0-255]\.[0-255]\.[0-255]$' <<< "$1"
    res=$?
    if [[ $res -ne 0 ]]; then  
        printf "%s\n" "$1 is not a valid IPv4 address"
    fi
}

cd(){
    local dir error
    while : ; do
        case $1 in
            --) break;;
            -*) shift;;
            *) break;;
        esac
     done
    dir=$1
    if [[ -n "$dir" ]]; then
        pushd "$dir"
    else
        pushd "$HOME"
    fi 2>/dev/null
    error=$?
    if [[ $error -ne 0 ]]; then
        builtin cd "$dir"
    fi
    return "$error"
} >/dev/null


pd(){
    popd
} >/dev/null

menu(){
    local item_number=0
    local selected_item=""
    for item in $@; do
        ((item_number++))
        printf "%d. %-30s" ${item_number} "${item}"
        if [[ $((item_number%2)) -eq 0 ]]; then
            printf "\n"
        fi
    done
    read -p "Select an item number: " item_number
    selected_item=$(eval printf "%s" "\$$item_number")
    printf "%s\n" "${selected_item}"
    #return "${selected_item}"
}

cdm(){

   local printed_columns=4
   local longest_string=0
   local current_string=0
   local num_lines=0
   local number_to_add=0
   local counter=0
   declare -a possible_locations
   local target_directory=${HOME}
   local term_width=$(tput cols)
   local left_adjustment=$(bc <<< "scale=0; ${term_width}/${printed_columns}*0.1")
   #local left_adjustment=$(bc <<< "scale=0; ${term_width}/${printed_columns}*0.9")
   local left_adjustment=${left_adjustment%.*}
   local num_directories=${#DIRSTACK[@]}
   local column_width=$(bc <<< "scale=0; ${term_width}/${printed_columns}")
   column_width=${column_width%.*}
   local target_remainder=$((printed_columns-1))

   while read line; do
       if [[ "$line" == "$PWD" ]]; then
           continue
       else
           case ${possible_locations[*]} in 
           *"${line}:"*) ;; 
           * ) 
               possible_locations+=("${line}:")
               ;;
           esac
       fi
       current_string=$(printf "%s" "${line}" | awk '{print length}')
       [ ${current_string} -gt ${longest_string} ] && longest_string=${current_string}
   done <<< $(dirs -l -p)

   if [[ ${longest_string} -gt ${left_adjustment} ]]; then
       left_adjustment=$((longest_string+5))
   fi

   while [ ${left_adjustment} -ge ${column_width} -a ${printed_columns} -gt 1 ] ; do
           ((printed_columns--))
           target_remainder=$((printed_columns-1))
           column_width=$(bc <<< "scale=0; ${term_width}/${printed_columns}")
           column_width=${column_width%.*}
   done

   display_number=1

   for position in $(seq 0 $((${#possible_locations[@]}-1))); do
       if [[ $((position%printed_columns)) -eq ${target_remainder} ]]; then
           printf "%-${left_adjustment}s\n" "${display_number}. ${possible_locations[${position}]%%:}" #| fold -w ${term_width}
       else 
           printf "%-${left_adjustment}s\t" "${display_number}. ${possible_locations[${position}]%%:}" #| fold -w ${term_width}
       fi
       ((display_number++))
   done

   printf "\n"
   read -p "Select a directory to move into: " index

   if [ -n ${index} ]; then
        case ${index} in
          q|Q) return 0;;
        esac
       if [ ${index} -le ${#possible_locations[@]} ]; then
           ((index--))
           builtin cd ${possible_locations[${index}]%%:}
       else 
           printf "%s\n" "${index} is not a valid choice"
           return 1
      fi
   fi

   unset target_directory num_directories num_lines num_lines pwd_occurances pwd_locations \
   selected_item number_to_add possible_locations
}

c(){
    printf "\033[0;0H\033[J"
}

        
#cdm(){
##    item_number=1
##    for directory in ${DIRSTACK[@]}; do
##
##       printf "%d. %s"
#   local target_directory=${HOME}
#   local num_directories=${#DIRSTACK[@]}
#   local num_lines=0
#   declare -a pwd_locations
#   local pwd_occurances=0
#   local number_to_add=0
#   local counter=0
#   while read line; do
#       directory=$(awk '{print $2}' <<< $line)
#       if [[ "$directory" == "$PWD" ]]; then
#           pwd_locations+=($counter)
#           ((pwd_occurances++))
#           #echo $pwd_occurances
#           ((counter++))
#           continue
#       else
#           ((num_lines++))
#           printf "%-30s \t" "$num_lines $directory"
#           if [ $((num_lines%3)) -eq 0 -a ${num_lines} -gt 0 ]; then
#               printf "\n"
#           fi
#       fi
#       ((counter++))
#   done <<< $(dirs -l -v -p)
#   printf "\n"
#   echo ${pwd_locations[@]}
#   read -p "Select a directory to move into: " index
#   if [ -n ${index} -a  ${index} -lt ${#DIRSTACK[@]} ]; then
#       #echo "valid index"
#       #echo $pwd_occurances
#        if [[ ${pwd_occurances} -gt 0 ]]; then
#            echo "occurances"
#            for location in $(seq 0 $((pwd_occurances-1))); do
#                if [[ ${pwd_locations[${location}]} -le ${index} ]]; then
#                    ((number_to_add++))
#                    #echo "added number"
#                fi
#            done
#        fi
#       ((index+=number_to_add))
#       #target_directory=${DIRSTACK[${index}]}
#       #cd $target_directory
#       echo $index
#   fi
#   unset target_directory num_directories num_lines num_lines pwd_occurances pwd_locations \
#   selected_item number_to_add
#}




kill-old-sessions(){
    my_tty=$(tty | cut -d'/' -f 3-)
    ps -ef | grep ssh | grep -v ${my_tty} | awk '{print $2}' | xargs -I {} kill -9 {}
    #sessions=$(ps -ef | grep ssh | grep -v ${my_tty} | awk '{print $2}')
   # printf "My tty is ${my_tty} and the other sessions are ${sessions}\n"
   # printf ${sessions} | xargs kill -9
}
set -o vi

#bind -m vi-insert '"\C-i": dynamic-complete-history'
#bind -m vi '"\C-i": dynamic-complete-history'


CDPATH="./:/home/jharmon:/etc:/var:"



change-colors() {
	
	color='31'

	case $1 in 
		
		34 | blue | Blue | BLUE)
			color='34'
			;;

		36 | cyan | Cyan | CYAN)
			color='36'
			;;

		32 | green | Green | GREEN)
			color='32'
			;;

		35 | purple | Purple | PURPLE)
			color='35'
			;;

		31 | red | Red | RED) 
			color='31'
			;;

		37 | white | White | WHITE)
			color='37'
			;;

		33 | yellow | Yellow | YELLOW)
			color='33'
			;;
			
		90 | grey | Gray | GREY)
			color='90'
			;;

		91 | bright_red)
			color='91'
			;;

		92 | bright_green)
			color='92'
			;;
		93 | bright_yellow)
			color='93'
			;;
		
		94 | bright_blue)
			color='94'
			;;

		95 | bright_magenta)
			color='95'
			;;

		96 | bright_cyan)
			color='96'
			;;

		97 | bright_white)
			color='97'
			;;

		--list | -l)
			printf "Possible colors:\n"
			printf "\033[00;31m31: red\033[00m\n"
			printf "\033[00;32m32: green\n"
			printf "\033[00;33m3: yellow\n"
			printf "\033[00;34m34: blue\n"
			printf "\033[00;35m35: purple\n" 
			printf "\033[00;36m36: cyan\n"
			printf "\033[00;37m37: white\n"
#			printf "\033[00;30m90: grey\n"
			printf "\033[00;91m91: bright_red\n"
			printf "\033[00;92m92: bright_green\n"
			printf "\033[00;93m93: bright_yellow\n"
			printf "\033[00;94m94: bright_blue\n"
			printf "\033[00;95m95: bright_magenta\n"
			printf "\033[00;96m96: bright_cyan\n"
			printf "\033[00;97m97: bright_white\n"
			;;

	esac

	close_format='\[\033[00m\]'
	start_format='\[\033[00;'$color'm\]'
	PS1="$start_format"'$USER@$HN:$PWD>'"$close_format"' '
	PS2="$start_format"'continue->'"$close_format"
}


alias vim="vimx -u ~/.vimrc"
set -o vi
bind '"jk":vi-movement-mode'
bind -m vi '"\e." insert-last-argument'
bind -m vi-insert '"\e." insert-last-argument'
bind -m vi-insert '"\C-x" edit-and-execute-command'
bind -m vi '"\C-x" edit-and-execute-command'


if [[ -f $/home/{SUDO_USER}/.vimrc ]]; then
    alias vim="vim -u /home/${SUDO_USER}/.vimrc"
fi


break_prompt(){
    PS1='\[\033[2J\]\[\033[3J\]\[\033[1;1H\]\[\033[08;31m\]$USER@$HOSTNAME:$PWD>\[\033[?25l;08;32m\] '
}

save_prompt(){
   PS1='          .     .  .      +     .      .          .
     .       .      .     #       .           .
        .      .         ###            .      .      .
      .      .   "#:. .:##"##:. .:#"  .      .
          .      . "####"###"####"  .
       .     "#:.    .:#"###"#:.    .:#"  .        .       .
  .             "#########"#########"        .        .
        .    "#:.  "####"###"####"  .:#"   .       .
     .     .  "#######""##"##""#######"                  .
                ."##"#####"#####"##"           .      .
    .   "#:. ...  .:##"###"###"##:.  ... .:#"     .
      .     "#######"##"#####"##"#######"      .     .
    .    .     "#####""#######""#####"    .      .
            .     "      000      "    .     .
       .         .   .   000     .        .       .
.. .. ..................O000O........................ ...... ...\[\033[?25h\]\[\033[01;31m\]$USER@$HN:$PWD>\[\033[00;32m\] '
}

long_prompt(){
    PS1='\[\033[01;31m\]$USER@$HOSTNAME:$PWD>\[\033[00;32m\] '
}

short_prompt(){
    PS1='\[\033[01;31m\]$USER@$HOSTNAME:\W>\[\033[00;32m\] '
}

long-prompt(){
    PS1='\[\033[00;31m\]$USER@$HOSTNAME:$PWD>\[\033[00;32m\] '
}

short-prompt(){
    PS1='\[\033[00;31m\]$USER@$HOSTNAME:\W>\[\033[00;32m\] '
}

set-selinux (){
options=$(getopt -o pm:n --long print,mode:,now -- "$@")
eval set -- "$options"
	mode="enforcing"
	print=""
	apply_now=""
	#do_return=""
	runtime_mode=$(getenforce)
	for arg; do
		case $arg in
			-p | --print)
			print="true"
				;;
			-m | --mode)
				shift
				mode="$1"
				if [[ $mode != "enforcing" ]] && [[ $mode != "permissive" ]] && [[ $mode != "disabled" ]]; then 
					printf "\033[38;05;204m\"${mode}\" is not an allowed selinux mode\n"
					echo "Allowed modes are: Enforcing, Permissive, or Disabled"
					return 1
				fi
				mode=$(tr [a-z] [A-Z] <<<${mode:0:1})${mode:1}
				;;
			-n | --now)
				apply_now="true"	
				;;
		esac
	done
	sed  -i.backup  "s/\(^SELINUX=\).*/\1$mode/1" /etc/selinux/config && printf "\033[38;5;027mPersistent mode changed to: ${mode}\033[00m\n"
	printf "\033[38;5;165mRuntime Mode is: ${runtime_mode}\n"
	[[ -n ${apply_now} ]] && setenforce ${mode}
	[[ -n ${print} ]] && printf "\033[38;5;220m###############SELinux Persistent Configuration###############\n" && cat /etc/selinux/config

}

show-pypath(){
for path in $(python3 -c 'import sys; print(sys.path)' | tr ',' '\n' | sed -n '2,$p' | sed "s/[]]//g" | sed  "s/'//g"); do echo $path; ls $path; printf "\n\n\n\n"; done
}

create-role(){
    cd ${ANSIBLE_PATH}/roles
    mkdir ${1}
    cd ${1}
    mkdir vars tasks templates files defaults handlers
}



# 415 
#asdfasfasf415afsdfasdfafs
#12393384158372019

# >>> conda initialize >>>
# !! Contents within this block are managed by 'conda init' !!
__conda_setup="$('/home/jharmon/anaconda3/bin/conda' 'shell.bash' 'hook' 2> /dev/null)"
if [ $? -eq 0 ]; then
    eval "$__conda_setup"
else
    if [ -f "/home/jharmon/anaconda3/etc/profile.d/conda.sh" ]; then
        . "/home/jharmon/anaconda3/etc/profile.d/conda.sh"
    else
        export PATH="/home/jharmon/anaconda3/bin:$PATH"
    fi
fi
unset __conda_setup
# <<< conda initialize <<<

