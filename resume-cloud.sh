#!/bin/zsh

export NO_PROXY=localhost,127.0.0.1
minikube start --container-runtime=docker --network-plugin=cni --cni=calico
minikube addons enable registry

lines=$(kubectl get svc -n services -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.spec.ports[].nodePort}{"\n"}{end}')
len=$(echo $lines | wc -l)
for ((i = 1; i <= $len; i++)); do
	name=$(echo $(echo $lines | awk 'NR==ctr' ctr="$i" | tr -d '\n' | awk '{ print $1 }'))
	port=$(echo $(echo $lines | awk 'NR==ctr' ctr="$i" | tr -d '\n' | awk '{ print $2 }'))
	kubectl port-forward $name -n services $port:3306 --address 192.168.1.65 &
done
