#!/bin/python3 

import os
from stat import *
from datetime import datetime

def get_metadata(file):
    metadata = dict()
    stat = os.stat(file)
    metadata['mode'] = oct(stat.st_mode)[-4:]
    metadata['uid'] = stat.st_uid
    metadata['gid'] = stat.st_gid
    metadata['mtime'] = datetime.fromtimestamp(stat.st_mtime).strftime("%Y-%m-%d %H:%M:%S")
    metadata['ctime'] = stat.st_ctime.fromtimestamp(stat.st_mtime).strftime("%Y-%m-%d %H:%M:%S")
    return metadata

def compare_metadata(md1, md2):
    diff_md = dict()
    for key in md1.keys():
        try:
            if (md1[key] != md2[key]):
                diff_md[key] = [md1[key], md2[key]]
        except KeyError as ke:
            pass
    if diff_md:
        return diff_md
    else:
        return False
            