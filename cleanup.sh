#!/usr/bin/env bash
# Clean up files older than $TIME days or so

TIME=31
WEBROOT=/var/www/paste

find $WEBROOT -type f -mtime +$TIME -name '*.paste' -execdir rm -- '{}' \;

