ifndef MASTER
MASTER=2
endif

start: create-config start-instances 

stop: stop-instances delete-config cleanup-volume-data

create-config:
	@echo "======================================================================="
	@echo "Creating the MYSQL configuration map"
	@echo "======================================================================="
	-kubectl create -f demo-cm.yaml

start-instances:
	@echo "======================================================================="
	@echo "Starting the MYSQL instances stateful set"
	@echo "======================================================================="
	-kubectl create -f demo-srv.yaml
	-kubectl create -f demo-ss.yaml

stop-instances:
	@echo "======================================================================="
	@echo "Stopping the MYSQL instances stateful set"
	@echo "======================================================================="	
	-kubectl delete -f demo-ss.yaml
	-kubectl delete -f demo-srv.yaml

delete-config:
	@echo "======================================================================="
	@echo "Deleting the MYSQL configuration map"
	@echo "======================================================================="	
	-kubectl delete -f demo-cm.yaml

cleanup-volume-data:
	if [ "$(KUBESLAVE0)" != "" ]; then ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i $(KUBE_SSH_KEY) opc@$(KUBESLAVE0) 'sudo rm -rf /mysql-data' ; fi
	if [ "$(KUBESLAVE1)" != "" ]; then ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i $(KUBE_SSH_KEY) opc@$(KUBESLAVE1) 'sudo rm -rf /mysql-data' ; fi
	if [ "$(KUBESLAVE2)" != "" ]; then ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i $(KUBE_SSH_KEY) opc@$(KUBESLAVE2) 'sudo rm -rf /mysql-data' ; fi
	if [ "$(KUBESLAVE3)" != "" ]; then ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i $(KUBE_SSH_KEY) opc@$(KUBESLAVE3) 'sudo rm -rf /mysql-data' ; fi

# =========================================================================== #
# Utility functions:
# =========================================================================== #

is-initialised:
	@echo "======================================================================="
	@echo "Have the init containers completed? (i.e. reached sleep 10000)"
	@echo "======================================================================="	
	kubectl logs mysql-0 -c mysql-init | tail -n 1
	kubectl logs mysql-1 -c mysql-init | tail -n 1
	kubectl logs mysql-2 -c mysql-init | tail -n 1

cluster-status:
	@echo "======================================================================="
	@echo "Cluster status"
	@echo "======================================================================="
	kubectl get statefulsets
	kubectl get pods -o wide
	kubectl exec -it mysql-0 -- mysqlsh --log-level=DEBUG3 --uri root:password@mysql-$(MASTER):3306 --py -e \
	    "import json; \
	     status = dba.get_cluster('devCluster').status(); \
	     parsed = json.loads(str(status)); \
         print json.dumps(parsed, indent=4, sort_keys=True)"

mysql-0:
	kubectl exec -it mysql-0 -- mysql -u root -ppassword
mysql-1:
	kubectl exec -it mysql-1 -- mysql -u root -ppassword
mysql-2:
	kubectl exec -it mysql-2 -- mysql -u root -ppassword

mysql-log-0:
	kubectl logs mysql-0
mysql-log-1:
	kubectl logs mysql-1
mysql-log-2:
	kubectl logs mysql-2

err-logs-0:
	kubectl exec -it mysql-0 -- cat /var/lib/mysql/mysqld.err
err-logs-1:
	kubectl exec -it mysql-1 -- cat /var/lib/mysql/mysqld.err
err-logs-2:
	kubectl exec -it mysql-2 -- cat /var/lib/mysql/mysqld.err

mysqlsh-log-0:
	kubectl exec -it mysql-0 -- cat /root/.mysqlsh/mysqlsh.log
mysqlsh-log-1:
	kubectl exec -it mysql-1 -- cat /root/.mysqlsh/mysqlsh.log
mysqlsh-log-2:
	kubectl exec -it mysql-2 -- cat /root/.mysqlsh/mysqlsh.log
