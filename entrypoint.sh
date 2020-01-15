#!/bin/bash
set -e

if [ "$1" = 'app' ]; then
	rm -rf /data/all/css
	cp -a css /data/all/css
	setfacl -m u:sftpuser:rwx /data/all
	find /data/all -type f -exec setfacl -m u:sftpuser:rw {} \;
	find /data/all -type d -exec setfacl -d -m u:sftpuser:rwx {} \;
	/etc/init.d/ssh start
    exec app "$@"
fi

exec "$@"
