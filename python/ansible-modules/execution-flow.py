import json
import yaml
import re
import os 
import sys

# Return the root path for the roles to be prepended to any include or import tasks
def set_role_root(file_path):
    return os.path.dirname(file_path)

def load_file(file_path):
    with open(file_path, 'r') as file:
        return file.read()
    
def load_yaml_file(file_path):
    yaml_struct = yaml.load(load_file(file_path))
    return yaml_struct

def parse_block(ansible_block): # parse one item in an ansible task list, for now only looking for certain keywords
    recursive_keywords = {
        'include_tasks': True, 
        'include_role': {
            'tasks_from': True
        },
        'import_tasks': True, 
        'import_role': {
            'tasks_from': True
        },
        'import_playbook': True 
    }
    
