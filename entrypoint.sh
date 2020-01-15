#!/bin/bash
set -e

if [ "$1" = 'app' ]; then
	if [ ! -d "/data/all/css" ]; then
		mkdir /data/all/css
	fi
	cp -a css/* /data/all/css/
	setfacl -m u:sftpuser:rwx /data/all
	find /data/all -type f -exec setfacl -m u:sftpuser:rw {} \;
	find /data/all -type d -exec setfacl -d -m u:sftpuser:rwx {} \;
	/etc/init.d/ssh start
    exec app "$@"
fi

exec "$@"
