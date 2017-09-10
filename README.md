##**goJam**  
Learning Go..  
For each project the path should be set as %GOPATH%/src/<your_project>  
 
  
#######################################################    
##**Expose the minikube docker env**  
```
minikube start
eval $(minikube docker-env --shell=bash)  
  
docker login    
docker pull store/oracle/weblogic:12.2.1.2    
docker rmi store/oracle/weblogic:12.2.1.2  
docker logs <container_id> | grep password  
[docker run -d store/oracle/weblogic:12.2.1.2]  
[docker run -d -p 7001:7001 store/oracle/weblogic:12.2.1.2]  
[docker run -d -p 7002:7001 store/oracle/weblogic:12.2.1.2]  
[docker logs <container_id> | grep password]  
  
#Run weblogic in minikube  
kubectl run weblogic --image=store/oracle/weblogic:12.2.1.2  
kubectl expose deployment weblogic --type=NodePort  
kubectl get pods  
kubectl get deployments    
kubectl get services  
minikube service weblogic --url  
kubectl delete service weblogic  
kubectl delete deployment weblogic  
```
#######################################################  
  

#######################################################  
##**Build and create docker image for our operator**  
```
kubectl exec -it my-pod --container main-app -- /bin/bash

Pre-req
----------
minikube 1.7
go 1.9
go get -u github.com/golang/dep/cmd/dep
glide-v0.12.3
GNU make

minikube start
eval $(minikube docker-env --shell=bash)
docker login

make clean
make vendor
make build
make image
make push

#Apply the operator
kubectl apply -f dist/weblogic-operator.yaml
kubectl -n weblogic-operator get pods

#Start Weblogic Servers using operator
kubectl apply -f examples/server.yaml
kubectl get weblogicservers

#Check the new services created
kubectl get services
minikube service adminserver --url
minikube service manageserver --url

#Delete Weblogic Server using operator
kubectl delete weblogicserver managedserver


#Cleanup
kubectl delete weblogicservers --all
kubectl delete -n weblogic-operator deployment weblogic-operator
kubectl delete ns weblogic-operator
```

