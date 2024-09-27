#!/bin/python3


# Simple script designed to read a bunch of shell a=b lines and translate that to a python dict, and output as yaml or json as needed
# eventually want to make this a lookup plugin or something like that for ansible
import sys, yaml, json
import collections
import re


def validate_assignment(line) -> bool:
		open_doubles = 0
		open_singles = 0
		is_valid = False
		for index, character in enumerate(line):
			#print(character)
			if character == "'": 
				#print('single boi')
				if open_doubles > 0:
					continue
				elif open_singles < 1:
					open_singles += 1
					print(f'Opening single quote begins.\nRemaining line: {line[index:]}')
				else:
					open_singles -= 1
					print(f'Closing single quote\nRemaining line: {line[index:]}')
				
			elif character == '"':
				#print('double boi')
				if open_singles > 0:
					continue
				elif open_doubles < 1:
					open_doubles += 1
					print(f'Opening double quote begins\nRemaining line: {line[index:]}')
				elif open_doubles > 0:
					open_doubles -= 1
					print(f'Closing double quote\nRemaining line: {line[index:]}')
			elif character == '=':
				if open_singles == 0 and open_doubles == 0:
					return True
		return False



def parse_line(line_number, line):
	is_assignment = False
	if not '=' in line:
		print(f'Line {line_number} contains no shell variable assignment')
		return False
	else:

		kv = line.split('=')
		key = kv[0]
		value = kv[1]
		return {key: value}
		
def parse_shell_file(shell_file):
	shell_data = dict()
	with open(shell_file, 'r') as sf:
		for line_number, line in enumerate(sf.readlines()):
			line_result = parse_line(line_number, line)
			shell_data.update(line_result)
	return shell_data

def main():
	if len(sys.argv) < 2:
		print("This script requires a file to be provided")
		exit(1)
#	shell_file = sys.argv[1]
#	shell_data = parse_shell_file(shell_file)
#	with open("./shell.yml", 'r') as sy:
		#for line in sy.readlines():
			#validate_assignment(line)
		##sy.write(yaml.dump(shell_data, default_flow_style=False))
	with open(sys.argv[1], 'r') as input_file:
		for line in input_file.readlines():
			print(validate_assignment(line))

if __name__ == '__main__':
	main()
