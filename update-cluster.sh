#!/bin/zsh

if [ -z "$1" ]; then
	echo "ERROR: enter name of cluster"
fi

if [ -z "$2" ]; then
	echo "ERROR: enter replica count"
fi

helm upgrade $1 mysql-chart --atomic --set replicaCount=$2

