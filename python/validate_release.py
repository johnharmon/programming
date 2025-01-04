#!/bin/python3
import os
import hashlib
import json 

def hash_file(path):
    sha = hashlib.sha256()
    with open(path, 'rb') as file:
        while True:
                buff = file.read(32000)
                if not buff:
                    break
                else:
                    sha.update(buff)
    return sha.hexdigest()

def generate_manifest():
    outer_dict = dict()
    with open('/tmp/released_manifest.json', 'w') as release_info_file:
        for root, dirname, filenames in os.walk(os.path.curdir ):
            for filename in filenames:
                absolute_path = os.path.join(root, filename)
                outer_dict[absolute_path] = dict()
                file_hash = hash_file(absolute_path)
                file_size = os.path.getsize(absolute_path)
                outer_dict[absolute_path]['Size'] = file_size
                outer_dict[absolute_path]['Hash'] = file_hash
        json_output = json.dump(outer_dict, release_info_file )

def compare_manifests():
    delivered_manifest = json.load(open('./released_manifest.json', 'r'))
    missing_files = dict()
    extra_files = dict()
    invalid_files = dict()
    valid_files = dict()
    for root, dirnames, filenames in os.walk(os.path.curdir):
        for filename in filenames:
           relative_path = os.path.join(root, filename)
           real_path = os.path.realpath(relative_path)
           try:
               delivered_manifest[relative_path]
           except (KeyError):
              extra_file_size = os.path.getsize(relative_path)
              extra_file_hash = hash_file(relative_path)
              extra_files[relative_path] = {
                  "Size": extra_file_size,
                  "Hash": extra_file_hash
              }
           else:
               local_hash = hash_file(relative_path)
               local_size = os.path.getsize(relative_path)
               if local_hash != delivered_manifest[relative_path]['Hash'] or local_size != delivered_manifest[relative_path]['Size']:
                   invalid_files[relative_path] = {
                       'Size': local_size,
                       'Hash': local_hash
                   }
                   delivered_manifest.pop(relative_path)
               else:
                   valid_files[relative_path] = {
                       'Size': local_size,
                       'Hash': local_hash
                   }
                   delivered_manifest.pop(relative_path)
    if len(delivered_manifest) > 0:
        missing_files = delivered_manifest
    
    print(f'There were {len(missing_files)} missing files, {len(invalid_files)} invalid files, and {len(extra_files)} extra files compared to the release manifest file')
    try: 
        with open('/tmp/valid_files.json', 'w') as output_file:
            json.dump(valid_files, output_file)
        with open('/tmp/invalid_files.json', 'w') as output_file:
            json.dump(invalid_files, output_file)
        with open('/tmp/missing_files.json', 'w') as output_file:
            json.dump(missing_files, output_file)
        with open('/tmp/extra_files.json', 'w') as output_file:
            json.dump(extra_files, output_file)
    except:
        pass

#if __name__ == '__main__':
#    main()

#def generate_new_manifest():
#    outer_dict = dict()
#    with open('/tmp/delivered_manifest.json', 'w') as release_info_file:
#        for root, dirname, filenames in os.walk(os.path.curdir ):
#            #outer_dict[f'{root}/{filename}'] = dict()
#            for filename in filenames:
#                absolute_path = os.path.join(root, filename)
#                outer_dict[absolute_path] = dict()
#                #print(absolute_path)
#                file_hash = hash_file(absolute_path)
#                file_size = os.path.getsize(absolute_path)
#                outer_dict[absolute_path]['Size'] = file_size
#                outer_dict[absolute_path]['Hash'] = file_hash
#                #release_info_file.writelines(f'{absolute_path}: {file_hash}\n')
#        #release_info_file.write(outer_dict)
#        #print(outer_dict)
#        json_output = json.dump(outer_dict, release_info_file )
#        #release_info_file.write(str(json_output))
#
