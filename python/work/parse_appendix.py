#!/bin/python3
import pandas as pd 
import re 
import sys 
import json

def read_nodes(csv_path, node_id=1):
    df = pd.read_csv(csv_path)
    matching_nodes = df.loc[df['Node_id'] == node_id]
    return matching_nodes

def filter_nodes(node_list):
    if len(node_list) == 1:
        return node_list
    else:
        return node_list.loc[node_list['Suite_type'] != 'playback']

def dump_node(node, json_path='~/ansible/roles/autogration/suite-forge/files/node-id.json'):
    node_dict = node.to_dict()
    with open(json_path, 'w') as f:
        try:
            json.dump(node_dict, f)
        except Exception as e:
            return False
        return True

def main():
    node_id = sys.argv[1]
    node_list = read_nodes('~/ansible/roles/autogration/suite-forge/files/nodes.csv', node_id)
    filtered_node_list = filter_nodes(node_list)
    result = dump_node(filtered_node_list)
    if result == False:
        exit(1)
    else:
        exit(0)

if __name__ == '__main__':
    main()