
eval "$(/opt/homebrew/bin/brew shellenv)"

alias ll='ls -l'
alias lst='ls -ltr'


function cdl (){
    cd $(ls -tr | tail -n 1)
}

eval "$(/opt/homebrew/bin/brew shellenv)"

#set -o vi
# source global definitions
#set -xv
#cd ~
# user specific environment
#path="$home/.local/bin:$home/bin:$path"
export path
export pythonpath="pythonpath:/home/jharmon/.local/lib/python3.9/site-packages"
export profile
export VIMRUNTIME=''
bindkey -e
bindkey '^[[1;5C' forward-word
bindkey '^[[1;5D' backward-word
#myvimrc='/home/jharmon/.vimrc'
#export myvimrc
ansible_path='/home/jharmon/programming/ansible/work'
# uncomment the following line if you don't like systemctl's auto-paging feature:
# export systemd_pager=

# user specific aliases and functions
alias ll='ls -ltr'
alias em='set +o vi; set -o emacs'
alias run9='podman run --rm -it registry.access.redhat.com/ubi9/ubi:latest /bin/sh'
alias cdg='cd ~/git-projects'
alias grc='vim ~/.config/ghostty/config'
alias vi='set +o emacs; set -o vi'
alias c='cd'
#alias vim='vim -u /home/jharmon/.vimrc'
alias lessi='less -i'
alias iless='less -i'
alias sm='ps -ef | grep -v grep | grep -e'
alias rm='rm -i'
alias cs='cd'
alias env-path='env | grep path | cut -d'=' -f 2 | tr ":" "\n"'
alias centos2='ssh 192.168.137.2'
alias ld='find . -maxdepth 1 -type d'
centos2='192.168.137.2'
export centos2
export histtimeformat="%h %d %h:%m:%s "
export histfilesize=1000000
export histsize=1000000
#shopt -s histappend
#prompt_command="$prompt_command; history -a"
export histcontrol=ignorespace
[[ ${sudo_user} == 'jharmon' ]] && touch /home/jharmon/.root_history && histfile='/home/jharmon/.root_history'
#function get-content {

#ifs=$'\n'

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
	#grep -e '[0-9]{1,2}/[^[:space:]].*' <<< $1 1>/dev/null && next_path=$(echo $1 | cut -d'/' -f2- | grep -v '^-')
	grep -e '[0-9]{1,2}/.*' <<< $1 1>/dev/null && next_path=$(echo $1 | cut -d'/' -f2- | grep -v '^-')
#	printf "after string creation\n"	
	for level in $(seq 1 $up_levels); do 
		up_string+='../'
	done
	result_string=${up_string}${next_path}
	cd ${result_string}
}
#ps1='\[\033[38;5;199m\]j\[\033[00;32m\]h\[\033[00;33m\]a\[\033[00;34m\]r\[\033[00;35m\]m\[\033[00;36m\]o\[\033[00;32m\]n\[\033[00;31m\]@$hostname:$pwd>\[\033[00;32m\] '
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
#				if [[ $1 == 'red' ]] || [[ $1 == 'red' ]] || [[ $1 == 'red' ]];
#				then
#					colorset='31'
#				elif [[ $1 == 'blue' ]] || [[ $1 == 'blue' ]] || [[ $1 == 'blue' ]];
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
#	ps1='$user@$hn:$pwd> '
#	ps2='continue->'
#
#
#}


show-colors (){ 
	options=$(getopt -o ls:t:o: --longoptions list,string:,trail:,stdout:,short,format -- "$@")
	eval set -- "$options"

	standard_output='color #'
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
		printf "\033[01;30m#30(black) \033[01;31m#31(red) \033[01;32m#32(green)  \033[01;33m#33(yellow) \033[01;34m#34(blue) \033[01;35m#35(magenta) \033[01;36m#36(cyan) \033[01;37m#37(white)\033[00m\n"
		unset short
		return

	elif [[ -n $format ]]; then
		 echo 'all color codes must be preceeded by \033[<attribute_code>;<background_code>;<foreground_code>m'
		 printf "(note: multiple attribute codes can be used, simply separate them via semicolons like the rest of the codes)\n"
		 printf "foreground color codes are:\n"
		printf "\033[01;30m#30(black) \033[01;31m#31(red) \033[01;32m#32(green)  \033[01;33m#33(yellow) \033[01;34m#34(blue) \033[01;35m#35(magenta) \033[01;36m#36(cyan) \033[01;37m#37(white)\033[00m\n"
		printf "background color codes are:\n"
		printf "\033[01;40;37m#40(black) \033[01;41;37m#41(red) \033[01;42;37m#42(green)  \033[01;43;37m#43(yellow) \033[01;44;37m#44(blue) \033[01;45;37m#45(magenta) \033[01;46;37m#46(cyan) \033[01;47;31m#47(white)\033[00m\n"
		printf "attribute modifiers are as follows:\n"
		printf "\033[00m#00 for normal\n\033[01m#01 for bold\n\033[04m#04 for underscore\n\033[05;41m#05 for blinking background\n\033[00m\033[07m#07 for reversed\n\033[08m#08 for concealed \033[00m\n"
		printf "the last one is #08 for concealed incase your terminal actually reads the codes properly\n"
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
    grep -e '^[0-255]\.[0-255]\.[0-255]\.[0-255]$' <<< "$1"
    res=$?
    if [[ $res -ne 0 ]]; then  
        printf "%s\n" "$1 is not a valid ipv4 address"
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
        pushd "$home"
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
    read -p "select an item number: " item_number
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
   local target_directory=${home}
   local term_width=$(tput cols)
   local left_adjustment=$(bc <<< "scale=0; ${term_width}/${printed_columns}*0.1")
   #local left_adjustment=$(bc <<< "scale=0; ${term_width}/${printed_columns}*0.9")
   local left_adjustment=${left_adjustment%.*}
   local num_directories=${#dirstack[@]}
   local column_width=$(bc <<< "scale=0; ${term_width}/${printed_columns}")
   column_width=${column_width%.*}
   local target_remainder=$((printed_columns-1))

   while read line; do
       if [[ "$line" == "$pwd" ]]; then
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
   read  "index?select a directory to move into: " 

   if [ -n ${index} ]; then
        case ${index} in
          q|q) return 0;;
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

#c(){
#    printf "\033[0;0h\033[j"
#}

        
#cdm(){
##    item_number=1
##    for directory in ${dirstack[@]}; do
##
##       printf "%d. %s"
#   local target_directory=${home}
#   local num_directories=${#dirstack[@]}
#   local num_lines=0
#   declare -a pwd_locations
#   local pwd_occurances=0
#   local number_to_add=0
#   local counter=0
#   while read line; do
#       directory=$(awk '{print $2}' <<< $line)
#       if [[ "$directory" == "$pwd" ]]; then
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
#   read -p "select a directory to move into: " index
#   if [ -n ${index} -a  ${index} -lt ${#dirstack[@]} ]; then
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
#       #target_directory=${dirstack[${index}]}
#       #cd $target_directory
#       echo $index
#   fi
#   unset target_directory num_directories num_lines num_lines pwd_occurances pwd_locations \
#   selected_item number_to_add
#}




kill-old-sessions(){
    my_tty=$(tty | cut -d'/' -f 3-)
    ps -ef | grep ssh | grep -v ${my_tty} | awk '{print $2}' | xargs -i {} kill -9 {}
    #sessions=$(ps -ef | grep ssh | grep -v ${my_tty} | awk '{print $2}')
   # printf "my tty is ${my_tty} and the other sessions are ${sessions}\n"
   # printf ${sessions} | xargs kill -9
}
#set -o vi

#bind -m vi-insert '"\c-i": dynamic-complete-history'
#bind -m vi '"\c-i": dynamic-complete-history'


cdpath="./:/home/jharmon:/etc:/var:"



change-colors() {
	
	color='31'

	case $1 in 
		
		34 | blue | blue | blue)
			color='34'
			;;

		36 | cyan | cyan | cyan)
			color='36'
			;;

		32 | green | green | green)
			color='32'
			;;

		35 | purple | purple | purple)
			color='35'
			;;

		31 | red | red | red) 
			color='31'
			;;

		37 | white | white | white)
			color='37'
			;;

		33 | yellow | yellow | yellow)
			color='33'
			;;
			
		90 | grey | gray | grey)
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
			printf "possible colors:\n"
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
	ps1="$start_format"'$user@$hn:$pwd>'"$close_format"' '
	ps2="$start_format"'continue->'"$close_format"
}


#alias vim="vimx -u ~/.vimrc"
#bind '"jk":vi-movement-mode'
#bind -m vi '"\e." insert-last-argument'
#bind -m vi-insert '"\e." insert-last-argument'
#bind -m vi-insert '"\c-x" edit-and-execute-command'
#bind -m vi '"\c-x" edit-and-execute-command'
#

if [[ -f $/home/{sudo_user}/.vimrc ]]; then
    alias vim="vim -u /home/${sudo_user}/.vimrc"
fi


break_prompt(){
    ps1='\[\033[2j\]\[\033[3j\]\[\033[1;1h\]\[\033[08;31m\]$user@$hostname:$pwd>\[\033[?25l;08;32m\] '
}

save_prompt(){
   ps1='          .     .  .      +     .      .          .
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
.. .. ..................o000o........................ ...... ...\[\033[?25h\]\[\033[01;31m\]$user@$hn:$pwd>\[\033[00;32m\] '
}

long_prompt(){
    ps1='\[\033[01;31m\]$user@$hostname:$pwd>\[\033[00;32m\] '
}

short_prompt(){
    ps1='\[\033[01;31m\]$user@$hostname:\w>\[\033[00;32m\] '
}

long-prompt(){
    ps1='\[\033[00;31m\]$user@$hostname:$pwd>\[\033[00;32m\] '
}

short-prompt(){
    ps1='\[\033[00;31m\]$user@$hostname:\w>\[\033[00;32m\] '
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
					echo "allowed modes are: enforcing, permissive, or disabled"
					return 1
				fi
				mode=$(tr [a-z] [a-z] <<<${mode:0:1})${mode:1}
				;;
			-n | --now)
				apply_now="true"	
				;;
		esac
	done
	sed  -i.backup  "s/\(^selinux=\).*/\1$mode/1" /etc/selinux/config && printf "\033[38;5;027mpersistent mode changed to: ${mode}\033[00m\n"
	printf "\033[38;5;165mruntime mode is: ${runtime_mode}\n"
	[[ -n ${apply_now} ]] && setenforce ${mode}
	[[ -n ${print} ]] && printf "\033[38;5;220m###############selinux persistent configuration###############\n" && cat /etc/selinux/config

}

show-pypath(){
for path in $(python3 -c 'import sys; print(sys.path)' | tr ',' '\n' | sed -n '2,$p' | sed "s/[]]//g" | sed  "s/'//g"); do echo $path; ls $path; printf "\n\n\n\n"; done
}

create-role(){
    cd ${ansible_path}/roles
    mkdir ${1}
    cd ${1}
    mkdir vars tasks templates files defaults handlers
}



# 415 
#asdfasfasf415afsdfasdfafs
#12393384158372019

# >>> conda initialize >>>
# !! contents within this block are managed by 'conda init' !!
if [ $? -eq 0 ]; then
    eval "$__conda_setup"
else
    if [ -f "/home/jharmon/anaconda3/etc/profile.d/conda.sh" ]; then
        . "/home/jharmon/anaconda3/etc/profile.d/conda.sh"
    else
        export path="/home/jharmon/anaconda3/bin:$path"
    fi
fi
unset __conda_setup
# <<< conda initialize <<<
alias cdc='cd ~/programming'

function go-daily() {
	project_name=$(printf $* | tr ' ' '-')
	dir_name="$(date +%y%m%d)-${project_name}"
	mkdir "${dir_name}"
	touch ${dir_name}/${project_name}.go ${dir_name}/go.mod
}
alias cdp='cd ~/git-projects/programming'

function cdl() {
	cd $(ls -tr | tail -n 1)
}


function python-daily() {
	project_name=$(printf $* | tr ' ' '-')
	dir_name="$(date +%y%m%d)-${project_name}"
	mkdir "${dir_name}"
	pythonfile=${dir_name}/${project_name}.py
	touch ${pythonfile}
	echo '' >> ${pythonfile}
	sed -i '1 i\#!/bin/python3' ${pythonfile}
	echo 'def main():' >> ${pythonfile}
	echo '	pass' >> ${pythonfile}
	echo '' >> ${pythonfile}
	echo -e "if __name__ == '__main__':\n\tmain()" >> ${pythonfile}
}
#export editor=/usr/bin/vim
export EDITOR=/opt/homebrew/bin/nvim
export GIT_EDITOR=/opt/homebrew/bin/nvim



function cdf (){ # function to cd into local directory based off .*<pattern>.* matching in current directory
	file_list=$(find . -maxdepth 1 -type d -iregex ".*${1}.*" 2>/dev/null | head -n 1)
	 

}

function gcmg (){
git clone https://fusionetixJohn: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@github.com/fusionetix/Music-Guessing-Game.git
}

export FUSIONETIX_ACCESS_TOKEN='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'
alias vim='nvim'
tmux





