# MySQL EE Operator

[![build status](https://gitlab-odx.oracle.com/odx/weblogic-operator/badges/master/build.svg)](https://gitlab-odx.oracle.com/odx/weblogic-operator/commits/master)

The Oracle MySQL [Operator][1] creates, configures, and
manages MySQL Enterprise clusters on Kubernetes.

**While fully usable, this is currently alpha software and should be treated as
such.  There may be backwards incompatible changes up until the first major
release.**

For support, please utilize the #mysql-operator-sup channel on Slack.  For bugs,
please open an [issue](https://gitlab-odx.oracle.com/odx/weblogic-operator/issues).

## Features

 * Create and delete MySQL clusters
 * Backup and restore of MySQL databases
 * Scheduled backups to Object Storage (OCI, S3 etc)

The MySQL Operator leverages Kubernetes [CustomResourceDefinitions][2] to define custom resources.

## Requirements

 * Kubernetes 1.7.0 +
 * (Optional) BMCS volume support

If you are running on Oracle OCI you will need to make sure your cluster supports persistent volumes via the
bmcs-flexvolume-driver and bmcs-volume-provisioner to make use of block storage volumes.

## Documentation

See [docs](docs) for detailed information about installing and configuring the mysql-operator.

[1]: https://coreos.com/blog/introducing-operators.html
[2]: https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/
