minikube delete
export NO_PROXY=localhost,127.0.0.1
minikube start --container-runtime=docker --network-plugin=cni --cni=calico
minikube addons enable registry
kubectl create namespace mysql
kubectl create namespace services
kubectl label namespaces mysql purpose=prod --overwrite=true
kubectl label namespaces services purpose=forwarding --overwrite=true
kubectl label namespaces kube-system purpose=k8s --overwrite=true
kubectl create -f network.yaml
