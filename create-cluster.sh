#!/bin/zsh

set -e

if [ -z "$1" ]; then
	echo "ERROR: enter name of cluster"
  exit 1;
fi

if [ -z "$2" ]; then
	echo "ERROR: enter replica count"
  exit 1;
fi

if [ -z "$3" ]; then
	echo "ERROR: enter mysql password"
  exit 1;
fi

helm install $1 mysql-chart --atomic --set replicaCount=$2 --set mysql.root_pass=$3
port=$(kubectl get svc $1 -n services -n services -o jsonpath='{.spec.ports[].nodePort}')
kubectl port-forward $1 -n services $port:3306 --address 192.168.1.65 1>/dev/null 2>&1 &
