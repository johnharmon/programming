#!/bin/python3
from collections import deque 
import sys
import os

def search_line(search_character, line):
    characters = dict()
    characters['('] = ['(', ')']
    characters['{'] = ['{', '}']
    characters['['] = ['[', ']']
    characters['<'] = ['<', '>']
    characters['"'] = ['"', '"']
    characters["'"] = ["'", "'"]
    search_stack = deque()
    for character in line:
        if character == characters[search_character][0]:
            search_stack.append(character)
        elif character == characters[search_character][1]:
            try:
                search_stack.pop()
            except Exception as e:
                return f'Unmatched closing {characters[search_character][1]}'
    if len(search_stack) > 0:
        return f'Unmatched opening {characters[search_character][0]}'
    else:
        return True

def search_file(search_character, file_path):
    wrong_lines = list()
    with open(file_path, 'r') as search_file:
         line_number = 1
         for line in search_file.readlines():
            result = search_line(search_character, line)
            if  result == True:
                pass
            else:
                wrong_lines.append(f'Line: {line_number}: {result}')
            line_number += 1

    return wrong_lines

def main():
    characaters = dict()
    characaters['('] = ['(', ')']
    characaters['{'] = ['{', '}']
    characaters['['] = ['[', ']']
    characaters['<'] = ['<', '>']
    characaters['"'] = ['"', '"']
    characaters["'"] = ["'", "'"]

    search_character = sys.argv[1]
    file_path = sys.argv[2]
    result = search_file(search_character, file_path)

    if len(result) > 0:
        for line in result:
            print(f'{line} in {os.path.relpath(file_path)}')
    else:
        print(f'This file does not have any unmatched {search_character} or {characaters[search_character][1]} on any line')
 
if __name__ == '__main__':
    main()




                