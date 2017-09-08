# MySQLOperator Developer Guide

---
## A. Configure k8s cluster for mysqloperator

To develop the mysqloperator and mysqlagent you will need access to a functioning kubernetes cluster with a defined set of namespaces and secrets configured. 

The following steps will initialise these resources and should be run against your vanilla k8s cluster.

### 1. Add mysql-operator namespace

```
kubectl create ns mysql-operator
```

This is required for the mysqloperator pod.

### 2. Configure k8s cluster docker registry 'odx-docker-pull-secret' secret.

```
kubectl -n mysql-operator create secret docker-registry odx-docker-pull-secret \
--docker-server="registry.oracledx.com" \
--docker-username="agent" \
--docker-password="XXX" \
--docker-email="k8s@oracle.com"

kubectl create secret docker-registry odx-docker-pull-secret \
--docker-server="registry.oracledx.com" \
--docker-username="agent" \
--docker-password="XXX" \
--docker-email="k8s@oracle.com"
```

These are required to allow k8s to access the mysqloperator and mysqlagent docker images created by the build process so that it can run them.


### 3. Configure k8s cluster 'mysql-root-user-secret' secret.

```
kubectl create secret generic --from-literal=password=mytestpass mysql-root-user-secret
```

This is required to assign the root user password to the mysql instances.


### 4. Configure k8s cluster 'bmc credentials' secret.

Create a credentials file containing the required passwords and keys. A template can be found [here](../examples/bmcs-backup-credentials.yaml).

```
provider: bmcs
bucket: backups
credentials:
  user: ${USER_OCID}
  fingerprint: ${USER_FINGERPRINT}
  key_file: |
    -----BEGIN RSA PRIVATE KEY-----
    ${BASE64_KEYDATA}
    -----END RSA PRIVATE KEY-----
  tenancy: ${TENANCY_OCID}
  region: ${BMC_REGION}
  namespace: ${BMC_NAMESAPCE}
```

... then upload it to the k8s cluster:

```
kubectl create secret generic --from-file ${PATH_TO_CREDENTIALS}
```

This is required so the the mysqlagent can access a bmc object storage service when performing 'backuup' and 'restore' operations.

---

## B. Build and run mysqloperator and mysqlagent

### 1. Build and push mysql-agent sidecar images

First, you need to build the mysqlagent image and export the version to environment:

```
make agent-push
export MYSQL_AGENT_VERSTION=tlangfor-20170830105344
```

This is required to perform backup and restore operations.

You then have several options for main mysqloperator image; you can build and run the golang binary locally:


### 2a. Build and run mysqloperator locally as a golang runtime

```
make run-dev
```

... or, run the controller from the image in k8s cluster:

### 2b. Build and run as a docker image (requires docker secrets)

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

## C. Create a MySQLCluster using the mysql operator

The simplest cluster required for devlopment consists of one node cluster with an associated persistent volume to mount the database. 

### 1. Create a cluster

A suitable template for this can be found [here](../examples/cluster-with-volume.yaml). This uses the names k8s secrets we have previously configured. Now create the cluster:

```
kubectl create -f examples/cluster-with-volume.yaml
```

### NB: Create k8s instance local mount directory:

If the cluster does not start and on inspection of the pod state you see a persisent volume releated mount error, then, you may need to log onto the host node to create the mount point:

```
ssh -o UserKnownHostsFile=/dev/null \
    -o StrictHostKeyChecking=no \
    -i /Users/tlangfor/.ssh/obmc-bristoldev/obmc-bristoldev \
    opc@129.146.43.204

mkdir /tmp/data
```

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

## C. Create a Backup of a MySQLCluster using the mysql operator


### 1. Create a Backup resource

A backup can be created for this cluster from the specification [here](../examples/backup.yaml), which has been configured to use the cluster and secrets defined in the previous steps. You can create the backup as follows:

```
kubectl create -f examples/backup.yaml
```

### 2. Investigate the backup

You can list and determine details of any backup. You can determine its current state (e.g. 'completed') and the location of stored backup image.

```
kubectl get mysqlbackups
kubectl describe mysqlbackup example-snapshot-backup
```

You can also check the log of the mysql-agent to see the output logs of the backup process:

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

