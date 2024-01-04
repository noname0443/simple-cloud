minikube delete
minikube start --container-runtime=docker --network-plugin=cni --cni=calico
kubectl create namespace mysql
kubectl label namespaces mysql purpose=prod --overwrite=true
kubectl label namespaces kube-system purpose=k8s --overwrite=true
kubectl accept -f network.yaml