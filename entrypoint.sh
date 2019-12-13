#!/bin/bash
set -e

if [ "$1" = 'app' ]; then
	find /data/all -type f-exec setfacl -m u:sftpuser:rw {} \;
	find /data/all -type d -exec setfacl -d -m u:sftpuser:rwx {} \;
    	exec app "$@"
fi

exec "$@"
