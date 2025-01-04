#!/bin/python3

import os
import sys
import re

def prarse_args():
    return_dict = dict()
    for index in range(len(sys.argv)):
        if not os.path.exists(sys.argv[0]):
            exit(0)
        else:
            target_file = sys.argv[0]
            return_dict['file'] = target_file
        if sys.argv[1]:
            return_dict['search_string'] = sys.argv[1]
    return return_dict
        
def create_node_groups(target_file):
    file_contents = open(target_file, 'r').readlines()
    node_ids = dict()
    for line in file_contents:
        node_id = re.search('(?<=_)[0-9]{1,3}(?=_)', line).group(0)
        if not node_id.group(0) in node_ids:
            node_ids[node_id] = dict()
        else:
            node_id_property= re.search('^_[0-9]{1,3}_[a-Z0-9_](?==)', line).group(0)
            node_id_property_value = line.split('=')[1]
            node_ids[node_id][node_id_property] = node_id_property_value
    return node_ids

def parse_groups(groups, search_string):
    for key in groups.keys():
        for property in key.keys():
            if groups[key][property] 

def main():
    results = parse_args()
    create_node_groups()

if __name__ == '__main__':
    main()
