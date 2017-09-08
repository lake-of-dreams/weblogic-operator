##**goJam**  
Learning Go..  
For each project the path should be set as %GOPATH%/src/<your_project>  
  
  
#######################################################  
##**Start minikube**  
```
https_proxy=www-proxy-idc.in.oracle.com:80 minikube start --docker-env HTTP_PROXY=www-proxy-idc.in.oracle.com:80 --docker-env HTTPS_PROXY=www-proxy-idc.in.oracle.com:80 --docker-env NO_PROXY=*.oraclecorp.com,*.oracle.com,192.168.99.0/24
```  

##**Run sample**  
```
kubectl run hello-minikube --image=gcr.io/google_containers/echoserver:1.4 --port=8080  
kubectl expose deployment hello-minikube --type=NodePort  
kubectl get pods  
kubectl get deployments  
kubectl get services  
minikube service hello-minikube --url  
kubectl delete service hello-minikube  
kubectl delete deployment hello-minikube  
minikube stop  
```
#######################################################  
  
  
#######################################################    
##**Expose the minikube docker env**  
```
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
##**Build docker image for our sample server.go and run in minikube** 
##This can be similar to weblogic running.##  
```
docker build -t gojamserver:v0 .
docker tag gojamserver:v0 fmwpltqa/gojamserver:v0
docker push fmwpltqa/gojamserver:v0

kubectl run hello-server --image=fmwpltqa/gojamserver:v0 --port=7777  
kubectl expose deployment hello-server --type=NodePort  
kubectl get deployments  
kubectl get services  
minikube service hello-server --url  
kubectl delete service hello-server  
kubectl delete deployment hello-server  
```
#######################################################  

#######################################################  
##**[WIP]Build docker image for our example operator**  
##Write an operator for the sample server to understand what an operator is##  
```
cd cmd
go get -v -d
docker build -t gojamoperator:v0  .
docker tag gojamoperator:v0 fmwpltqa/gojamoperator:v0
docker push fmwpltqa/gojamoperator:v0

kubectl run hello-server --image=fmwpltqa/gojamoperator:latest --port=9999
kubectl expose deployment jamserver-operator --type=NodePort  
kubectl get pod  
kubectl get services  
minikube service hello-server --url  
kubectl delete service hello-server  
kubectl delete deployment hello-server  
```
#######################################################


```
kubectl create -f manifests/gojam-server-crd    #Create CustomResourceRefinition
/apis/gojam.com/v1/namespaces/*/jamservers/
kubectl apply -f specs/gojam-server             #Deploy the server as the CustomType defined above
kubectl get jamserver



kubectl exec -it my-pod --container main-app -- /bin/bash

```

