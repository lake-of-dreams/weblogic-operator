This repository contains a Proof of Concept for a Kubernetes Operator for WebLogic.

Please refer to the project <a href="https://gitlab-odx.oracle.com/marnelso/weblogic-operator/wikis/home">wiki</a> for more information. 


**Pre-requisites**  
```
minikube 1.7
git                                             #Git (with bash) for Windows v2.14 
                                                (support many linux commands)
go 1.9
dep                                             #Dependency Management 
                                                (go get -u github.com/golang/dep/cmd/dep)
make                                            #MinGW or GNU Make for Windows
```

**Start minikube and export docker env** 
```
minikube start
eval $(minikube docker-env --shell=bash)
docker login
```

**Build _weblogic-operator_**
```
#Clone the source into $GOPATH/src/weblogic-operator

make clean
make vendor                                     #Uses dep to populate vendors
make build
``` 

**Create _weblogic-operator_ image and push** 
```
make image
make push                                       #Pushes to docker.io/fmwplt/weblogic-operator
``` 

**Create Secret to be used for pulling _WebLogicManagedServer_ image from registry**
```
kubectl create secret docker-registry weblogic-docker-store --docker-server=docker.io \
--docker-username=YOUR_USERNAME --docker-password=YOUR_PASSWORD --docker-email=YOUR_EMAIL
``` 

**Create CRD of type _WebLogicManagedServer_ into k8s**
```
kubectl apply -f manifests/weblogic-crd.yaml    #Creates custom object of type WebLogicManagedServer
``` 

**Deploy _weblogic-operator_ into k8s**
```
kubectl apply -f manifests/weblogic-operator.yaml
kubectl -n weblogic-operator get pods
``` 

**Create objects of type _WebLogicManagedServer_**
```
kubectl apply -f examples/server.yaml
kubectl get weblogicmanagedservers,services
``` 

**Delete objects of type _WebLogicManagedServer_**
```
kubectl delete weblogicmanagedserver managedserver
``` 

**Cleanup**
```
kubectl delete weblogicmanagedservers --all
kubectl delete -n weblogic-operator deployment weblogic-operator
kubectl delete ns weblogic-operator
```


