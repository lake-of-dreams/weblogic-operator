##**Pre-requisites**  
```
minikube 1.7
go 1.9
dep                 #Dependency Management (go get -u github.com/golang/dep/cmd/dep)
```

##**Start minikube and export docker env** 
```
minikube start
eval $(minikube docker-env --shell=bash)
docker login
```

##**Build _weblogic-operator_**
```
make clean
make vendor         #Uses dep to populate vendors
make build          #Build binary and files to dist/
``` 

##**Create _weblogic-operator_ image and push** 
```
make image
make push           #Pushes to docker hub
``` 

##**Apply CRD for _weblogic-operator_ into k8s**
```
kubectl apply -f dist/weblogic-crd.yaml          #Creates custom object of type WeblogicServer
``` 

##**Deploy the _weblogic-operator_ into k8s**
```
kubectl apply -f dist/weblogic-operator.yaml
kubectl -n weblogic-operator get pods
``` 

##**Create objects of type _WeblogicServer_**
```
kubectl apply -f examples/server.yaml
kubectl get weblogicservers,services
``` 

##**Create objects of type _WeblogicServer_**
```
kubectl delete weblogicserver managedserver
``` 

##**Cleanup**
```
kubectl delete weblogicservers --all
kubectl delete -n weblogic-operator deployment weblogic-operator
kubectl delete ns weblogic-operator
```