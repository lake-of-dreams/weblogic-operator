# MySQL InnoDB Cluster Investigations

## Within Docker

To bring up a 3 node cluster within docker, with 1 master (R/W) and 2 slaves (R/O):

```
make -f Makefile.docker start
```

To delete the cluster:

```
make -f Makefile.docker stop
```

Reference: https://github.com/mattlord/Docker-InnoDB-Cluster

## Within Kubernetes

Note: For this proof-of-concept to work, you'll need a cluster with 3 slaves. This is due to 
the use of host-path volumes, and our assumption that each of the MYSQL instances will run on a different 
slave, each with their own host-path volume directory.

First run the following script to install the docker registry secret within Kubernetes:

```
./generate-docker-registry-secret.sh
```

To bring up a 3 node cluster using a stateful set within Kubernetes, with 1 master (R/W) and 2 slaves (R/O):

```
make start
```

To delete the cluster:

```
make stop
```

## Tested Scenarios

1. **Cluster creation - WORKS** Created cluster with 3 nodes (primary + 2 secondaries). Cluster comes up OK with all nodes online.

2. **Secondary failure - WORKS** Delete (using kubectl) one of the secondary instance pods. Kubernetes will then restart this pod, and the init
   container kicks in an does the cluster rejoin.

3. **Primary failure - WORKS** Delete (using kubectl) the primary instance pod. The cluster will failover and make one of the secondaries
   the new primary. Kubernetes will restart the old primary pod, and the init container will ensure that this rejoins the 
   cluster as a new read-only secondary.

4. **Server container dies - WORKS** Log into the relevant Kubernetes node and kill just the MYSQL server container by doing a docker stop.
   Kubernetes will ensure that this container gets restarted. The init container will then notice that the server has been restarted
   and wake up and configure it to rejoin the cluster.

5. **Scale out - WORKS** Update the number of replicas from 3 to 4, then do a kubectl apply. The now secondary
   will get started and join the cluster automatically. Note: You'll need to ensure you have at least as many
   Kubernetes nodes as replicas, in order for each to have its own host based volume.

6. **Scale in - WORKS** Update the number of replicas from 4 to 3, then do a kubectl apply. The secondary is killed.

7. **Rolling upgrade - NOT WORKING** Updating the image currently breaks the cluster. This is due to the readiness 
   checks not being implemented, thus when doing the rolling upgrade there is a time when all the servers are down and 
   initialising. Therefore, when the first one comes back up it assumes that there is a cluster to join, but this
   is not the case. We need to implement the readiness probes such that each instance comes up fully before the next
   replica is upgraded. However, see below for some notes on investigations in this area. 

## Known Issues

* When scaling back in, the terminated instance is still known to the cluster and is marked as missing in the 
  cluster status output. Maybe we should try to remove this (using remove_instance? in mysqlsh)

* If a primary failovur has occurred prior to doing a scale out, the scale out fails because the 
  new instance assumes mysql-0 is the primary, and is is hardwired to login to this instance using mysqlsh.
  We need work out which one the primary is in this case.  

* We need to ensure that the resource constraints are set to be very low for the init containers, as 
  we dont want a larger than nessasary default resource limit on the init container constraining where 
  the pod can be scheduled.

## General Notes

### Use of readiness probes

Currently the stateful set doesn't use any liveness/readiness probes, and all (initial) the 
MYSQL instances come up together, the secondaries employing a wait loop to ensure that the 
cluster is created prior to attempting to join. Another way of doing this would be 
to use the readiness probe to (fully) bring up one instance at a time. However, I tried this
and it gave errors when the secondaries attempted to join the cluster, with the logs showing 
connection/network issues. I didn't fully understand the reason for this though, so this might need 
further investigation if we decide to go down this route.

### Manage the cluster from within the pods of within the cluster?

One question is whether to implement the cluster operations (i.e. create, add instance, rejoin instance)
inside the stateful set itself, maybe as a sidecar container in each replica pod, or handle
this externally in the controller. I spent some time trying to do the former. However,
this proves fiddly for a couple of reasons:

If the MYSQL server container crashes and gets restarted, then the sidecar init container
will also need to be rerun in order to bring the server back into the cluster. 
Not sure how to achieve this.

The logic is fiddly because from within the container, you dont have the overview. In general,
you need to use the mysqlsh to connect to one of the cluster instances then something like the following:

```
if cluster exists:
    if I was a memner but am now missing:
        rejoin the cluster [dba.get_cluster().rejoin_instance()]

    else:
        add myself to the cluster [dba.get_cluster().add_instance()]
else:
    create the cluster [dba.create_cluster()]
```

However, how do we know which instance to connect to? i.e. the sidecar container may
need to check all 3 to determine if this is the first time the cluster is created, or
if this pod is being revived after a failure. Doing this inline in the stateful set definition
gets messy. Also, we don't want to build our own image, as this make the upgrade
path more difficult, i.e. our image would have to keep up.

So, all in all, I think the option of orchestrating the cluster from the controller might be the best option...

### Cluster configuration, built in my.cnf or use mysqlsh to configure?

Currently the configuration is specified in a my.cnf file, which is set in Kubernetes 
using a config map. However, we might be able to get away with no configuration up front,
and use mysqlsh to configure the instances using dba.check_instance_configuration and 
dba.configure_local_instance. This might be better as it make the upgrade path more robust, i.e.
we wouldn't need to ensure that our configuration file is compatible with future MYSQL releases.

Looked into this further. The mysqlsh can be used to configure an instance for clustering. However, after
configuring (which updates /etc/my.cnf) the instance needs restarting, which is a problem as we are running in 
a container and a restart will lose the config! 

Also investigated taking the minimal config file (produced from mysqlsh) and using this directly. However,
this is also producing connection/network errors, i.e. the MYSQL instances cant talk to each other.
