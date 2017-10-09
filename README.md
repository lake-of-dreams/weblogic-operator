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
make                                            #MinGW Make for Windows
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

**Create Secret to be used for pulling _WebLogic_ image from registry**
```
kubectl create secret docker-registry weblogic-docker-store --docker-server=docker.io \
--docker-username=YOUR_USERNAME --docker-password=YOUR_PASSWORD --docker-email=YOUR_EMAIL
``` 

**Create CRD's of type _WebLogicDomain_ and _WebLogicManagedServer_ into k8s**
```
kubectl apply -f manifests/weblogic-crd.yaml    #Creates custom object of type WebLogicDomain and WebLogicManagedServer
``` 

**Configure Persistant Volume Storage**
```
# Currently uses hostPath type. Host folder used is /scratch
# In Windows/VirtualBox, modify the minikube machine to add a new shared volume to any location in host and
specify it to be mount as /scratch.  
  
kubectl apply -f manifests/persistant-volume.yaml
``` 

**Deploy _weblogic-operator_ into k8s**
```
kubectl apply -f manifests/weblogic-operator.yaml
kubectl -n weblogic-operator get pods
``` 

**Create objects of type _WebLogicDomain_**
```
#Domain will be created in persistant volume with managed servers named as managedserver-0...n and starts AdminServer
#Default credentials used are weblogic/welcome1.  
  
kubectl apply -f examples/domain.yaml
kubectl get weblogicdomains,services
``` 

**Create objects of type _WebLogicManagedServer_**
```
#Domain created in persistant volume will be used
#Next available server is calculated from $DOMAIN_HOME/serverList.json and starts the requested no:of servers
#Default credentials used are weblogic/welcome1.  
  
kubectl apply -f examples/server.yaml
kubectl get weblogicservers,services
``` 

**Delete objects of type _WebLogicDomain_**
```
kubectl delete weblogicdomain basedomain
``` 

**Delete objects of type _WebLogicManagedServer_**
```
kubectl delete weblogicmanagedserver basedomain-managedserver
``` 

**Cleanup**
```
kubectl delete weblogicmanagedservers --all
kubectl delete weblogicdomains --all
kubectl delete deployment weblogic-operator
```


