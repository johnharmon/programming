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
    return set(result)

def search_file(file: str):
    # use set object for holding only unique paths
    unique_paths = set()
    if os.path.isfile(file):
        with open(file, 'r') as input_file:
            for line in input_file.readlines():
                # Update set with union of new files read from each line
                unique_paths = unique_paths.union(search_line(line))
    iteration_set = unique_paths.copy()
    # Remove values from set if they are the current file name, would cause useless and infinite recursion
    for path in iteration_set:
        if path == file:
            unique_paths.discard(path)
    if unique_paths:
        # Convert it to a list before returning cause sets don't JSON
        return list(unique_paths)

def recursive_search(file: str, depth: int, calling_file = None):
    # Check to see if we are at our first level of recursion, directory changing will be dependent on calling file so relative paths remain intaact
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
        # Essentially test to see if the calling file is in a different path and change to there
        dname = os.path.dirname(calling_file)
        if os.path.exists(dname):
            os.chdir(os.path.dirname(calling_file))
    file_results = dict()
    file_results[file] = dict()
    if get_file_type(file):
        file_results[file]['Paths'] = search_file(file)
        # Limit recursive depth to 3 levels
        if file_results[file]['Paths'] and depth < 3:
            # Increment depth counter
            depth += 1
            for filepath in file_results[file]['Paths']:
                # Update this level's dictionary with the results of recursively searching the paths found in this file
                file_results[file].update(recursive_search(filepath, depth, calling_file = file))
    return file_results 

def main():
    file = sys.argv[1]
    results = recursive_search(file, depth = 0)
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


