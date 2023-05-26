#!/bin/python3
import sys
import os
import yaml
import re

node_file_location = os.path.join(os.path.dirname(sys.argv[0]), 'suite_properties.txt')

with open(node_file_location, 'r') as node_file:
    
    if sys.argv[1]:
        pattern = sys.argv[1]
        #print(pattern)
        script_path = os.path.dirname(sys.argv[0])
        lines = node_file.read().splitlines()
        iteration_list = lines.copy()
        property_regex = ''
        for idx in range(0,len(iteration_list)):
            if iteration_list[idx] == '':
                lines.remove(iteration_list[idx])

        #print(lines)
        with open(f'{script_path}/suite_attributes.txt', 'r') as suite_attributes:
            suite_attributes = suite_attributes.read().splitlines()
            property_regex = '('
            for idx in range(0,len(suite_attributes)):
                if suite_attributes[idx] == '':
                    continue
                if idx != len(suite_attributes)-2:
                    property_regex += f'{suite_attributes[idx]}|'
                else:
                    property_regex += f'{suite_attributes[idx]})'
            #print(property_regex)

        matching_lines = list()
        matching_nodes = dict()
        if re.search('^[0-9]{1,3}$', pattern):
            #print('node id search activated')
            matching_nodes[f'Node ID - {pattern}'] = list()
            for line in lines:
                #if re.search(f'^_{pattern}_.*{property_regex}.*$', line):
                if re.search(f'^_{pattern}_.*$', line):
                    if re.search(f'.*{property_regex}.*', line):
                        matching_nodes[f'Node ID - {pattern}'].append('Property: ' + line.strip("'"))

        else:
            #print('string pattern search activated')
            for line in lines:
                node_id = re.search('^(?:_)[0-9]{1,3}(?:_)', line)
                if node_id:
                    node_id = node_id.group().strip('_')
                #node_id = node_id.strip('_')
                    #print(f'.*{property_regex}.*{pattern}.*', line)
                    my_match =  re.search(f'.*{property_regex}.*{pattern}.*', line)
                    if my_match:
                        #print(my_match.start(), my_match.endpos)
                        if not f'Node ID - {node_id}' in matching_nodes:
                            matching_nodes[f'Node ID - {node_id}'] = list()
                        matching_nodes[f'Node ID - {node_id}'].append(f'Property: {line}')

        if matching_nodes:
            print(yaml.dump(matching_nodes, default_flow_style=False))
            #print(yaml.dump(lines, default_flow_style=False))