#!/bin/python3
import os
import re
import yaml
import copy

file_names = dict()

try:
    with open('./ansible.modules', 'r') as module_file:
        ansible_modules = module_file.readlines()
        for index in range(0,len(ansible_modules)):
            ansible_modules[index] = ansible_modules[index].split()[0]

except:
    print('unable to locate ansible.modules file in working directory!')
    exit(1)

for dirpath, dirnames, filenames in os.walk('./'):
    for filename in filenames:
        if re.match('.*\.yml$', filename) and not re.match('(vault.yml|vars.yml)', filename):
            file_info = dict()
            file_info['Possible_lines'] = list()
            with open(os.path.join(dirpath, filename), 'r') as yml_file:
                file_contents = yml_file.readlines()
                for line in file_contents:
                    my_match = re.match('^[\-\s]*[a-zA-Z\.\-_]+:[\s]*$', line)
                    if my_match:
                        if not re.match('^[\-\s]*(roles|tasks|collections|name|when|block|loop|vars|include|register):[\s]*$', line):
                            string_match = re.match('[a-zA-Z\.\-_]+', my_match.string)
                            #if string_match:
                            #    print(string_match.string)
                            line = line.strip().strip(':').strip('-').strip()
                            file_info['Possible_lines'].append(line)
            if len(file_info['Possible_lines']) > 0:
                file_names[os.path.join(os.path.realpath(dirpath), filename)] = file_info
                
with open('./ansible_possible_modules', 'w') as output_file:
    yaml.dump(file_names, output_file, default_flow_style=False)

file_names2 = copy.deepcopy(file_names)

for file in file_names.keys():
    for module in file_names[file]['Possible_lines']:
        if module in ansible_modules:
            while True:
                try:
                    file_names2[file]['Possible_lines'].remove(module)
                except:
                    break
    if len(file_names2[file]['Possible_lines']) == 0 :
        file_names2[file].pop('Possible_lines')
    if len(file_names2[file]) == 0:
            file_names2.pop(file)

with open('./ansible_missing_modules', 'w') as missing_module_file:
    yaml.dump(file_names2, stream=missing_module_file, default_flow_style=False)
