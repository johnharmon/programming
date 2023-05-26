#!/bin/python3
from genericpath import isfile
import os
import hashlib
import json
import sys
import re
import time

def parse_args():
    manifest_location = '.'
    write_location = None
    search_directory = None
    generate_manifest_bool = False
    compare_manifest_bool = False
    if len(sys.argv) > 1:
        for index in range(1, len(sys.argv)):
            if re.match('--manifest|-m', sys.argv[index]):
                try:
                    manifest_location = sys.argv[index+1]
                except:
                    print('An argument is required for the --manifest or -m parameter!')
                    exit(1)
                else:
                    try:
                        manifest_location = find_manifest(manifest_location)
                    except:
                        print(f'No manifest found at {manifest_location}!')
                        exit(1)
            if re.match('--directory|-d', sys.argv[index]):
                try:
                    search_directory = sys.argv[index+1]
                except:
                    print('An argument is required for the --directory or -d parameter!')
                    exit(1)
            if re.match('--generate|-g', sys.argv[index]):
                generate_manifest_bool = True
            if re.match('--compare|-c', sys.argv[index]):
                compare_manifest_bool = True
            if re.match('--write|-w', sys.argv[index]):
                try:
                    write_location = sys.argv[index+1]
                    if not os.path.isdir(write_location):
                        raise ValueError
                except:
                    print(f'A target directory must be given to the {sys.argv[index]} parameter!')
                    exit(1)
    return manifest_location, search_directory, generate_manifest_bool, compare_manifest_bool, write_location

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

def generate_manifest(search_directory=None):
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
    print(os.path.join(os.path.realpath(target_directory), 'release-manifest.json'))
    with open(os.path.join(target_directory, 'release-manifest.json'), 'w') as release_info_file:
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


def compare_manifests(manifest_location=None, search_directory=None, write_location = None):
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
                open(find_manifest(manifest_location), 'r'))
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

def main():
    manifest_location, search_directory, generate_manifest_bool, compare_manifest_bool, write_location = parse_args()
    #if manifest_location =
    if generate_manifest_bool:
        print('generating manifest')
        generate_manifest(search_directory)
    if compare_manifest_bool:
        compare_manifests(manifest_location=manifest_location, search_directory=search_directory, write_location = write_location)

if __name__ == '__main__':
    main()
