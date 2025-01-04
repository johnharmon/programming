#!/bin/python3
import sys
import os
import yaml
import re

def node_search(patterns):

    def integer_search(node_content, pattern_list, property_regex):
        # Create dictionary to store the data from the node that matches the integer given
        matching_nodes = dict()
        # Formatting dictionary key to be desired output, pobably should change this
        matching_nodes[f'Node ID - {pattern_list[0]}'] = list()
        # Integer search ignores all patterns after the integer
        pattern = pattern_list[0]
        for line in node_content:
            # Checks to see if the line contains the node id we are looking for
            if re.search(f'^_{pattern}_.*$', line):
                # Checks to see if the line for a matching node id is one of the important values we care about
                if re.search(f'.*{property_regex}.*', line, re.I):
                    # Append new line containing useful value to the list contained within the dictionary
                    matching_nodes[f'Node ID - {pattern}'].append('Property: ' + line.strip("'"))
        return matching_nodes

    def check_node(node, pattern_list, property_regex):
        iteration_list = pattern_list.copy()
        # Create copy of pattern_list that can have values removed from it without effecting the iteration process
        removal_list = pattern_list.copy()
        important_lines = list()
        matched = False
        for line in node:
            # Check if line contains one of the important properties loaded from suite_attributes.txt, append to important_lines if it does
            if re.search(f'.*{property_regex}.*', line, re.I):
                important_lines.append(line)
            # For each pattern in the list, check to see if the line matches all of them
            if not matched:
                for pattern in iteration_list:
                    if re.search(f'.*{pattern}.*', line, re.I):
                        # If the line matches the current pattern being checked, remove that pattern from the list
                        removal_list.remove(pattern)
                        # if all patterns have been matched, set matched to true and break the loop. We do not return here because we do not know if we have found all the important lines
                        if len(removal_list) == 0:
                            matched = True
                            break
            # Set the remaining patterns to be mached to the list of unmatched patterns after processing the previous line
            iteration_list = removal_list
        if matched:
            return important_lines 
        else:
            return None

    def string_search(node_content, pattern_list, property_regex):
        nodes = dict()
        matching_nodes = dict()
        # Iterate over all lines in the suite.properties file, splitting up the nodes into separate dictionaries of lists which contain all the lines relevant to a particular node
        # This approach allows nodes to be properly evaluated even if the lines were completely out of order and nodes values weren't grouped together
        for line in node_content:
            # Check to see if we could pull a valid integer node id from the line
            node_id = re.search('^(?:_)[0-9]{1,3}(?:_)', line)
            if node_id:
                # Strip out underscores from matched integer
                node_id = node_id.group().strip('_')
                node_id = f'Node ID - {node_id}'
                #Check to see if this node id already exists as a dictionary key, if not, then create it and make the value an empty list
                if not node_id in nodes:
                    nodes[node_id]=list()
                # Append the new line to the now existing node id dictionary entry
                nodes[node_id].append('Property: ' + line.strip("'"))
        # Pass off the series of nodes to the check_node function to see if they match the values we are looking for
        for node in nodes.keys():
            result = check_node(nodes[node], pattern_list, property_regex)
            if result:
                matching_nodes[node] = result 
        return matching_nodes
            
    node_file_location = os.path.join(os.path.dirname(patterns[0]), 'suite_properties.txt')
    with open(node_file_location, 'r') as node_file:
        
        if patterns:
            #pattern = sys.argv[1]
            # A list of all cli arguments given to the script, if they are not integers lines will be searched to see if they match these arguments when evaluating matching nodes
            pattern_list = patterns[1:]
            #print(pattern)
            script_path = os.path.dirname(patterns[0])
            lines = node_file.read().splitlines()
            real_lines = list()
            property_regex = ''
            matching_nodes = dict()
            matching_lines = list()
            # Filter out empty lines so that they don't have to get reprocessed during regex matching
            for idx in range(0,len(lines)):
                if lines[idx] != '':
                    real_lines.append(lines[idx])

        # Used to load a list of important properties for nodes we will want to print, instead of printing all properties for matching nodes
        # Loads these properties from a text file, right now hard coded as suite_attributes
        with open(f'{script_path}/suite_attributes.txt', 'r') as suite_attributes:
            suite_attributes = suite_attributes.read().splitlines()
            # The property regex variable is used to essentially generate a series of regex or expressions so we can pull out the lines that are useful to us without looping over the poossible useful values each time
            property_regex = '('
            for idx in range(0,len(suite_attributes)):
                if suite_attributes[idx] == '':
                    continue
                if idx != len(suite_attributes)-2:
                    property_regex += f'{suite_attributes[idx]}|'
                else:
                    property_regex += f'{suite_attributes[idx]})'

            # Checks to see if the first argument given is an integer, if so, it searches the suite.properties file for said node and returns it
            if re.search('^[0-9]{1,3}$', pattern_list[0]):
                return integer_search(node_content = real_lines, pattern_list = pattern_list, property_regex = property_regex)
            # If the first argument is a string, it starts a string search attempting to match the rest of the given patterns to specific lines
            else:
                return string_search(node_content = real_lines, pattern_list = pattern_list, property_regex = property_regex)

def main():
    matching_nodes = node_search(sys.argv)
    if matching_nodes:
        print(yaml.dump(matching_nodes, default_flow_style=False))
    else:
        print('\033[00;33mNo nodes matched the pattern(s) given!\033[00m')

if __name__ == '__main__':
    main()

