#!/bin/zsh

set -e

export NO_PROXY=localhost,127.0.0.1
minikube start --container-runtime=docker --network-plugin=cni --cni=calico
minikube cp ./resolv.conf /etc/resolv.conf
minikube addons enable registry
