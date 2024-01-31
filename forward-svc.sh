#!/bin/zsh

set -e

if [ -z "$1" ]; then
	echo "ERROR: enter name of cluster"
	exit 1
fi

name=$1
port=$(kubectl get svc -n services $name -o jsonpath='{.spec.ports[].nodePort}')

kubectl port-forward $name -n services $port:3306 --address 0.0.0.0 1>/dev/null 2>&1 &
