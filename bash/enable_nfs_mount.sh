#!/bin/bash

setup_files="/tmp/jharmon/bash-scripts/setup-files"

[[ -f "${setup_files}/nfs-automount.service" ]] && cp "${setup_files}/nfs-automount.service" /etc/systemd/system
[[ -f "${setup_files}/nfs_automount.sh" ]] && cp "${setup_files}/nfs_automount.sh" /usr/local/bin

chmod u+x /usr/local/bin/nfs_automount.sh
systemctl enable nfs-automount.service


