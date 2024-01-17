#!/bin/zsh

set -e

export NO_PROXY=localhost,127.0.0.1
minikube start --container-runtime=docker --network-plugin=cni --cni=calico
minikube cp ./resolv.conf /etc/resolv.conf
minikube addons enable registry

kubectl wait pod --for=condition=Ready --all -n services --timeout 600s

lines=$(kubectl get svc -n services -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.ports[].nodePort}{"\n"}{end}')
len=$(echo $lines | wc -l)
for ((i = 1; i <= $len; i++)); do
	name=$(echo $(echo $lines | awk 'NR==ctr' ctr="$i" | tr -d '\n' | awk '{ print $1 }'))
	port=$(echo $(echo $lines | awk 'NR==ctr' ctr="$i" | tr -d '\n' | awk '{ print $2 }'))
	kubectl port-forward $name -n services $port:3306 --address 192.168.1.65 1>/dev/null 2>&1 &
done
