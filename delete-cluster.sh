#!/bin/bash

name=$1
count=$2
pass=$3

if [ -z "$name" ]; then
	echo "ERROR: enter name of cluster"
	exit 1
fi

helm uninstall $name
kubectl delete pvc -n mysql -l app.kubernetes.io/instance=$1
