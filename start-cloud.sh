minikube delete
minikube start --container-runtime=docker --network-plugin=cni --cni=calico
kubectl create namespace mysql
kubectl label namespaces mysql purpose=prod --overwrite=true
kubectl accept -f network.yaml