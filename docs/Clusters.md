# Clusters

MySQL cluster examples.

### Create a cluster with a custom "MYSQL_ROOT_PASSWORD"

Create your own secret with a password field

```
$ kubectl create secret generic mysql-root-user-secret --from-literal=password=foobar
```

Create your cluster and reference it

```yaml
apiVersion: "mysql.oracle.com/v1"
kind: MySQLCluster
metadata:
  name: example-mysql-cluster-custom-secret
spec:
  replicas: 1
  secretRef:
    name: mysql-root-user-secret
```

### Create a cluster with a persistent volume

The following example will create a MySQL Cluster with a persistent local volume.

```yaml
---
apiVersion: v1
kind: PersistentVolume
metadata:
  labels:
    type: local
  name: mysql-local-volume
spec:
  accessModes:
  - ReadWriteMany
  capacity:
    storage: 10Gi
  hostPath:
    path: /tmp/data
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: manual
---
apiVersion: "mysql.oracle.com/v1"
kind: MySQLCluster
metadata:
  name: example-mysql-cluster-with-volume
spec:
  replicas: 1
  volumeClaimTemplate:
    metadata:
      name: data
    spec:
      storageClassName: manual
      accessModes:
        - ReadWriteMany
      resources:
        requests:
          storage: 1Gi
```
