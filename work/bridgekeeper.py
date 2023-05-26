#!/bin/python3
from genericpath import isfile
import os
import hashlib
import json
import sys
import re
import time
import socket
import psycopg2

def parse_args():
    manifest_location = '.'
    write_location = None
    search_directory = None
    generate_manifest_bool = False
    compare_manifest_bool = False
    manifest1 = None
    manifest2 = None
    output_file = None
    compare_file = None
    if len(sys.argv) > 1:
        for index in range(1, len(sys.argv)):
            if re.match('--file|-f', sys.argv[index]):
                if compare_manifest_bool == True:
                    print('\033[01;35mYou cannot use the -f and -c flag at the same time!\033[00m')
                    exit(1)
                try:
                    manifest_location = sys.argv[index+1]
                    compare_manifest_bool = True
                    if not os.path.isfile(manifest_location):
                        print(f'{manifest_location} is not a file on the filesystem!')
                        exit(1)
                except IndexError:
                    print('An argument is required for the --file or -f parameter!')
                    exit(1)
                else:
                    try:
                        manifest_location = find_manifest(manifest_location)
                    except:
                        print(f'No manifest found at {manifest_location}!')
                        exit(1)
            if re.match('--name|-n', sys.argv[index]):
                try:
                    generate_manifest_bool = True
                    if generate_manifest_bool:
                        print('You do not need to use the -g flag when using the -n flag')
                    output_file = sys.argv[index+1]
                except IndexError:
                    print(f'An argument must be provided to the --name or -n flag')
                    exit(1)
#                    cur_time = time.localtime()
#                    output_file = f'/tmp/{socket.gethostname()}-{cur_time.tm_year}-{cur_time.tm_mon:02d}-{cur_time.tm_mday:02d}.{cur_time.tm_hour:02d}:{cur_time.tm_min:02d}:{cur_time.tm_sec:02d}'
            if re.match('--directory|-d', sys.argv[index]):
                try:
                    search_directory = sys.argv[index+1]
                except:
                    print('An argument is required for the --directory or -d parameter!')
                    exit(1)
            if re.match('--generate|-g', sys.argv[index]):
                if generate_manifest_bool:
                        print('You do not need to use the -g flag when using the -n flag')
                generate_manifest_bool = True
            if re.match('--compare|-c', sys.argv[index]):
                if compare_manifest_bool == True:
                    print('\033[01;35mYou cannot use the -f and -c flag at the same time!\033[00m')
                    exit(1)
                else:
                    compare_manifest_bool = True
                try:
                    if os.path.isfile(sys.argv[index+1]):
                        if os.path.isfile(sys.argv[index+2]):
                            manifest1 = sys.argv[index+1]
                            manifest2 = sys.argv[index+2]
                            #print(manifest1)
                            #print(manifest2)
                        else:
                            print(f'{sys.argv[index+2]} is not an existing file!')
                            exit(1)
                    else:
                        print(f'{sys.argv[index+1]} is not an existing file!')
                        exit(1)
                except Exception as e:
                    #print(e)
                   # print(sys.argv)
                   # print(index)
                    print('did not recieve 2 files to compare')
                    exit(1)

               # try:
               #     compare_file = sys.argv[index+1]
               # except IndexError:
               #     print('An argument must be provided to the -c or --compare parameter!')
               #     exit(1)
            if re.match('--write|-w', sys.argv[index]):
                try:
                    write_location = sys.argv[index+1]
                    if not os.path.isdir(write_location):
                        raise ValueError
                except:
                    print(f'A target directory must be given to the {sys.argv[index]} parameter!')
                    exit(1)
    return manifest_location, search_directory, generate_manifest_bool, compare_manifest_bool, write_location, output_file, manifest1, manifest2

def find_manifest(path):
    if os.path.isdir(path):
        target = os.path.join(path, 'release-manifest.json')
    else:
        target = path
    if os.path.isfile(target):
        return target
    else:
        raise FileNotFoundError

def hash_file(path):
    sha = hashlib.sha256()
    try:
        with open(path, 'rb') as file:
            while True:
                buff = file.read(32000)
                if not buff:
                    break
                else:
                    sha.update(buff)
        return sha.hexdigest()
    except:
        return False

def get_total_size(path):
    #total_size = os.popen(f'du -b --summarize {path}').read().split()[0]
    total_size = 0
    for dirpath, dirnames, filenames in os.walk(path):
        for filename in filenames:
            try:
                total_size += os.path.getsize(os.path.join(dirpath, filename))
            except:
                pass
    return int(total_size)

def calculate_size_from_manifest(path):
    try:
        manifest_location = find_manifest(path)
    except:
        return get_total_size(path)
    else:
        with open(manifest_location, 'r') as manifest_file:
            total_size = 0
            manifest_entries = json.load(manifest_file)
            for key in manifest_entries.keys():
                total_size += int(manifest_entries[key]['Size'])
            return total_size

def generate_manifest(search_directory=None, output_file = None):
    checked_size = 0
    total_size = 0
    if search_directory:
        if os.path.isdir(search_directory):
            target_directory = search_directory
        else:
            print(f'Target directory of {search_directory} does not exist!')
    else:
        target_directory = os.path.curdir
    os.chdir(target_directory)
    target_directory = os.path.curdir
    total_size = get_total_size(target_directory)
    outer_dict = dict()
    #print(os.path.join(os.path.realpath(target_directory), 'release-manifest.json'))
    print(os.path.join(os.path.realpath('/tmp'), output_file))
    with open(output_file, 'w') as release_info_file:
        files_searched = 0
        for root, dirname, filenames in os.walk(target_directory):
            for filename in filenames:
                if os.path.join(root, filename) == release_info_file.name:
                    continue
                absolute_path = os.path.join(root, filename)
                outer_dict[absolute_path] = dict()
                file_hash = hash_file(absolute_path)
                if file_hash:
                    files_searched += 1
                    file_size = os.path.getsize(absolute_path)
                    checked_size += file_size
                    if files_searched == 3000:
                        files_searched = 0
                        print(
                            f'Last file checked: {os.path.realpath(absolute_path)}')
                        print(
                            f'Percent Completed: {int(checked_size/total_size*100)}%')
                    outer_dict[absolute_path]['Size'] = file_size
                    outer_dict[absolute_path]['Hash'] = file_hash
                else:
                    continue
        json.dump(outer_dict, release_info_file)

def compare_manifest(manifest_location=None, search_directory=None, write_location = None):
    if search_directory:
        if os.path.isdir(search_directory):
            target_directory = search_directory
        else:
            print(f'Target directory of {search_directory} not found!')
            exit(1)
    else:
        target_directory = os.path.curdir
    print(target_directory)
    os.chdir(target_directory)
    target_directory = os.path.curdir
    if not manifest_location:
        try:
            delivered_manifest = json.load(
                open('./release-manifest.json', 'r'))
        except:
            print('There is no manifest file inside of this directory!')
            exit(1)
    else:
        try:
            delivered_manifest = json.load(
                open(manifest_location, 'r'))
        except:
            print(f'No manifest file found at: {manifest_location}!')
            exit(1)
    #print(os.path.realpath)
    total_size = calculate_size_from_manifest(target_directory)
    original_size = total_size
    missing_files = dict()
    extra_files = dict()
    invalid_files = dict()
    valid_files = dict()
    files_checked = 0
    for root, dirnames, filenames in os.walk(target_directory):
        for filename in filenames:
            relative_path = os.path.join(root, filename)
            if relative_path == './release-manifest.json':
                #print(os.path.join(root,filename))
                continue
            real_path = os.path.realpath(relative_path)
            try:
                delivered_manifest[relative_path]
            except (KeyError):
                extra_file_size = os.path.getsize(relative_path)
                extra_file_hash = hash_file(relative_path)
                extra_files[relative_path] = {
                    'Size': extra_file_size,
                    'Hash': extra_file_hash
                }
            else:
                local_hash = hash_file(relative_path)
                if local_hash:
                    local_size = os.path.getsize(relative_path)
                    if local_hash != delivered_manifest[relative_path]['Hash'] or local_size != delivered_manifest[relative_path]['Size']:
                        invalid_files[relative_path] = {
                            'Size': local_size,
                            'Hash': local_hash
                        }
                        total_size -= local_size
                        delivered_manifest.pop(relative_path)
                    else:
                        valid_files[relative_path] = {
                            'Size': local_size,
                            'Hash': local_hash
                        }
                        total_size -= local_size
                        delivered_manifest.pop(relative_path)
                    files_checked += 1
                    if files_checked >= 1000:
                        files_checked = 0
                        print(f'Last file checked {os.path.realpath(relative_path)}')
                        print(f'Percent Completed: {int((original_size-total_size)/original_size*100)}%')

    if len(delivered_manifest) > 0:
        missing_files = delivered_manifest
    if write_location == None:
        print_location = 'None'
    else:
        print_location = write_location
    print(f'There were {len(missing_files)} missing files, {len(invalid_files)} invalid files, and {len(extra_files)} extra files compared to the release manifest.\nJSON files for each of these lists is available in {print_location}')
    if write_location:
        try:
            #os.makedirs('/tmp/css/release-manifest.jsons', exist_ok = True)
            with open(os.path.join(write_location, 'valid_files.json'), 'w') as output_file:
                json.dump(valid_files, output_file)
        except:
            print(f'Failed to write json file to {output_file}')
        try:
            with open(os.path.join(write_location, 'invalid_files.json'), 'w') as output_file:
                json.dump(invalid_files, output_file)
        except:
            print(f'Failed to write json to {output_file}')
        try:
            with open(os.path.join(write_location, 'missing_files.json'), 'w') as output_file:
                json.dump(missing_files, output_file)
        except:
            print(f'Failed to write json to {output_file}')
        try:
            with open(os.path.join(write_location, 'extra_files.json'), 'w') as output_file:
                json.dump(extra_files, output_file)
        except:
            print(f'Failed to write json to {output_file}')

def compare_manifests(manifest1=None, manifest2=None, write_location = None ):
    name1 = manifest1
    name2 = manifest2
    manifest1 = json.load(open(manifest1, 'r'))
    manifest2 = json.load(open(manifest2, 'r'))
    manifest1_keys = list(manifest1.keys())
    manifest2_keys = list(manifest2.keys())
    missing_files = dict()
    extra_files = dict()
    invalid_files = dict()
    valid_files = dict()
    #print(manifest1)
    #print(manifest2)

    for key in manifest1_keys:
        try: 
            manifest2[key]
        except KeyError:
            missing_files[key] = manifest1[key]
            manifest1.pop(key)
        else:
            if (manifest1[key]['Size'] != manifest2[key]['Size']) or manifest1[key]['Hash'] != manifest2[key]['Hash']:
                invalid_files[key] = manifest1[key]
                manifest2.pop(key)
    valid_files = manifest1
    extra_files = manifest2 

    print(f'There were {len(missing_files)} missing files, {len(invalid_files)} invalid files, and {len(extra_files)} extra files in {name2} compared to {name1}.\nJSON files for each of these lists is available in {write_location}')
    if write_location:
        try:
            #os.makedirs('/tmp/css/release-manifest.jsons', exist_ok = True)
            with open(os.path.join(write_location, 'valid_files.json'), 'w') as output_file:
                json.dump(valid_files, output_file)
        except:
            print(f'Failed to write json file to {output_file}')
        try:
            with open(os.path.join(write_location, 'invalid_files.json'), 'w') as output_file:
                json.dump(invalid_files, output_file)
        except:
            print(f'Failed to write json to {output_file}')
        try:
            with open(os.path.join(write_location, 'missing_files.json'), 'w') as output_file:
                json.dump(missing_files, output_file)
        except:
            print(f'Failed to write json to {output_file}')
        try:
            with open(os.path.join(write_location, 'extra_files.json'), 'w') as output_file:
                json.dump(extra_files, output_file)
        except:
            print(f'Failed to write json to {output_file}')

def write_to_db(outer_dict):
    connection = psycopg2.connect(
        database = 'tmp',
        user = 'tmp',
        password = 'tmp',
        host = 'tmp',
        port = 'tmp'
    )
    coneciton.autocommit = True
    cursor = connection.cursor()
    pgsql = '''CREATE TABLE BRIDGEKEEPER(
        Suite INT,
        Server varchar(50),
        baseline varchar(50),
        json TEXT, 
        PRIMARY KEY (suite, server, baseline)
        );
        '''
    cursor.execute(sql)
    sql = '''INSERT INTO BRIDGEKEEPER(Suite, Server, baseline, json) VALUES{};'''.format(suite, server, baseline, f'{outer_dict}')
    cursor.execute(sql)

def main():
    #print('main function invoked')
    manifest_location, search_directory, generate_manifest_bool, compare_manifest_bool, write_location, output_file, manifest1, manifest2 = parse_args()
    if not output_file:
        cur_time = time.localtime()
        output_file = f'/tmp/{socket.gethostname()}_{cur_time.tm_year}-{cur_time.tm_mon:02d}-{cur_time.tm_mday:02d}.{cur_time.tm_hour:02d}:{cur_time.tm_min:02d}:{cur_time.tm_sec:02d}.json'
    #if manifest_location =
    if generate_manifest_bool:
        print('generating manifest')
        generate_manifest(search_directory, output_file)
    if compare_manifest_bool:
        #print('comparison invoked')
        if manifest1 and manifest2:
#            print('comparing 2 manifest')
#            print('comparing manifests')
            compare_manifests(manifest1 = manifest1, manifest2 = manifest2, write_location = write_location)
        else:
            compare_manifest(manifest_location=manifest_location, search_directory=search_directory, write_location = write_location)

if __name__ == '__main__':
    main()
