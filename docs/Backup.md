# Backups

## Introduction

The MySQL Operator allows for on-demand and scheduled backups to be created.
On-demand backups can be created by submitting a MySQLBackup custom resource
while scheduled backups must be declared at the cluster level.

Whilst we plan to offer different options for backups, we currently support
only the [Oracle mysqlbackup][1] tool. This requires a commercial license.

### Types of backup

We support two kinds of backup policy.

 1. Snapshot (A full snapshot of a database)
 2. Delta (TBC)

### Credentials

All backups require object storage credentials in order for the mysql-agent to
persist backups. An example for Oracle BMCS is given below.

```yaml
provider: bmcs
bucket: backups
credentials:
  user: ocid1.user.oc1..aaaaaaaai77mql2xerv7cn6wu3nhxang3y4jk56vo5bn5l5lysl34avnui3q
  fingerprint: 8c:bf:17:7b:5f:e0:7d:13:75:11:d6:39:0d:e2:84:74
  key_file: |
    -----BEGIN RSA PRIVATE KEY-----
    <snip>
    -----END RSA PRIVATE KEY-----
  tenancy: ocid1.tenancy.oc1..aaaaaaaatyn7scrtwtqedvgrxgr2xunzeo6uanvyhzxqblctwkrpisvke4kq
  region: us-phoenix-1
  namespace: my-project
```

Now create a secret with the contents of the above yaml file.

```bash
$ kubectl create secret generic bmcs-upload-credentials --from-file=./examples/bmcs-backup-credentials.yaml
```

## On-demand backups

#### Snapshot

You can request a backup at any time by submitting a MySQLBackup CRD to the
operator. The secretRef is the name of a secret that contains your Object
Storage credentials.

```yaml
apiVersion: "mysql.oracle.com/v1"
kind: MySQLBackup
metadata:
  name: example-backup
cluster:
  name: mycluster
secretRef:
  name: bmcs-upload-credentials
```

#### Delta

TBC

## Scheduled backups

TBA (work in progress)

[1]: https://dev.mysql.com/doc/mysql-enterprise-backup/4.1/en/
