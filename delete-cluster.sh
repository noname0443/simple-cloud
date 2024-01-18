#!/bin/bash

if [ -z "$1" ]; then
	echo "ERROR: enter name of cluster"
	exit 1
fi

helm uninstall $1
kubectl delete pvc -n mysql -l app.kubernetes.io/instance=$1
