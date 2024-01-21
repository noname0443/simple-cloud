#!/bin/zsh
set -e

if [ -z "$1" ]; then
	echo "ERROR: enter name of cluster"
	exit 1
fi

if [ -z "$2" ]; then
	echo "ERROR: enter admin password of cluster"
	exit 1
fi

if [ -z "$3" ]; then
	echo "ERROR: enter replica count"
	exit 1
fi

if [ -z "$4" ]; then
	echo "ERROR: enter user creds 'john:pass'"
	exit 1
fi

curl -v -X POST "0.0.0.0:8080/api/update" -d "{\"name\": \"$1\", \"password\": \"$2\", \"replicas\": \"$3\"}" -H "Content-Type: application/json" --user $4
