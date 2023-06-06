#!/bin/python3

import re
import os
import sys
import typing
import yaml
import json

# Use file command to ensure file is human readable ascii text
def get_file_type(filepath: str):
    file_results = os.popen(f"file {filepath}").readlines()
    file_type = re.search('ASCII text', file_results[0])
    if file_type:
        return file_type
    else:
        return False

def search_line(line: str):
    # Squeeze whitespace
    line = re.sub( '\s+', ' ', line)
    line = line.strip()
    # Ensure leading space in line, fix for regex having trouble matching line start anchor or other characters
    line = ' ' + line
    # Regex with non-capturing look behind to see if it could be a filepath 
    file_start = r'(?<=[^A-z0-9])[./A-z0-9]+[A-z0-9/_.-]*'
    file_regex = file_start
    # Get list of all matching strings in the line
    result = re.findall(file_regex, line)
    # Create copy of list to iterate over
    iteration_list = result.copy()
    # remove elements if they don't exist on the system
    for entry in iteration_list:
        if not(os.path.exists(entry)):
            result.remove(entry)
            #print(f'removed {entry}')
    #print(result)
    return set(result)

def search_file(file: str):
    #filepaths = list()
    filepaths = list()
    unique_paths = set()
    if os.path.isfile(file):
        with open(file, 'r') as input_file:
            for line in input_file.readlines():
                unique_paths = unique_paths.union(search_line(line))
    iteration_set = unique_paths.copy()
    for path in iteration_set:
        if path == file:
            unique_paths.discard(path)
    if unique_paths:
        return list(unique_paths)

def recursive_search(file: str, depth: int, calling_file = None):
    #print(file)
    if not calling_file:
        if file[0] != '/':
            calling_directory = os.getcwd()
            target_path = os.path.dirname(os.path.abspath(os.path.join(calling_directory, file)))
            #print(target_path)
            os.chdir(target_path)
        else:
            #os.chdir(''.join(file.split('/')[0:-1]))
            if os.path.exists(file) and os.path.isfile(file):
                os.chdir(os.path.dirname(file))
            elif os.path.exists(file) and os.path.isdir(file):
                os.chdir(file)
    else:
        dname = os.path.dirname(calling_file)
        if os.path.exists(dname):
            os.chdir(os.path.dirname(calling_file))
    file_results = dict()
    file_results[file] = dict()
    if get_file_type(file):
        file_results[file]['Paths'] = search_file(file)
        #print(file_results)
        if file_results[file]['Paths'] and depth < 3:
            depth += 1
            for filepath in file_results[file]['Paths']:
                file_results[file][filepath] = dict()
                #file_results[file][filepath]['Paths'] = recursive_search(filepath, depth)
                file_results[file].update(recursive_search(filepath, depth, calling_file = file))
    return file_results 

def main():
    print(os.getcwd())
    file = sys.argv[1]
    results = recursive_search(file, depth = 0)
    #print(results)
    print(json.dumps(results, indent = 4))

if __name__ == '__main__':
    main()


    



'''    
    #file_start = '(?=^|[^A-z0-9/.])[./A-z0-9][A-z0-9/_.-]'
    #file_start = r'(?<=[^A-z0-9])[./A-z0-9]+[A-z0-9/_.-]*'
    #file_start = r'(?:[^A-z0-9/.!]|[\s])[./A-z0-9][A-z0-9/_.-]*'
    #file_start = '(?=^|[A-z])[A-z]'    #[./A-z0-9][A-z0-9/_.-]'
    #file_start = '[./A-z0-9][A-z0-9/_.-]'
    #file_start = re.compile(file_start)
    #file_middle = '[A-z0-9/_.-]'
    #file_middle = re.compile(file_middle)
'''


