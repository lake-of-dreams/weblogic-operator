# MySQLOperator Developer Guide


---

## B. Build and run mysqloperator and mysqlagent

### 2b. Build and run as a docker image ()

```
make deploy 
```
or
```
make push
make start
```
or
```
make push
kubectl apply -f dist/${OPERATOR_NAME}.yaml
```

We should now have a configured k8s environment for: creating clusters, taking backups, running e2e tests, etc.

---

### 3. Investigate the cluster

You should now be able to investigate the cluster:

```
kubectl get mysqlclusters
kubectl describe pod example-mysql-cluster-with-volume-0
```

... and log into a working mysql instance:

```
kubectl exec -it example-mysql-cluster-with-volume-0 -- bash -c 'mysql -uroot -p${MYSQL_ROOT_PASSWORD}'
```

... get the mysql container logs:

```
kubectl logs -f example-mysql-cluster-with-volume-0 -c mysql
```

... get the mysqlagent container logs:

```
kubectl logs -f example-mysql-cluster-with-volume-0 -c mysql-agent
```

---


## X. Tidy up resources

```
kubectl delete backups --all
kubectl delete mysqlclusters --all
kubectl delete -n mysql-operator deployment mysql-operator
```

