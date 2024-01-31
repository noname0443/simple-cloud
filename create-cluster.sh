#!/bin/zsh

set -e

if [ -z "$1" ]; then
	echo "ERROR: enter name of cluster"
  exit 1;
fi

if [ -z "$2" ]; then
	echo "ERROR: enter replica count"
  exit 2;
fi

if [ -z "$3" ]; then
	echo "ERROR: enter mysql user name"
  exit 3;
fi

if [ -z "$4" ]; then
	echo "ERROR: enter mysql password"
  exit 4;
fi

if [ -z "$5" ]; then
	echo "ERROR: enter mysql database name"
  exit 5;
fi

helm install $1 mysql-chart --atomic --set replicaCount=$2 --set mysql.user_password=$3 --set mysql.user=$4 --set mysql.user_database=$5 --timeout 10m0s
port=$(kubectl get svc $1 -n services -n services -o jsonpath='{.spec.ports[].nodePort}')
kubectl port-forward $1 -n services $port:3306 --address 0.0.0.0 1>/dev/null 2>&1 &
