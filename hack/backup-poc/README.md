# backup-poc

A dockerfile and some bash scripts for testing 'mysqlbackup' backup/restore operations.

TODO: Currently, hangs on takin snapshot after a restore
TODO: Does not implement storage service based archiving.

---

## example usage

### start and shell into docker instance

```
cd {project_home}/weblogic-operator/hack/backup-poc
# run 'clear' if the test container is already running...
make dev-docker-start
make dev-docker-shell
```

### source 'backup-poc.sh' to add debug and anlysis function to current environmet

```
source /backup-poc.sh
```

### init the 'mysqlbackup user' (if not already initalised)

```
/backup-poc.sh -u root -p root init
```

### create a new backup (primary snapshot backup)

```
/backup-poc.sh -u mysqlbackup -p new-password create_backup
```

### take a new delta (incremetal backup)

```
/backup-poc.sh -u mysqlbackup -p new-password create_delta
```

### restore from current backup image and associated delta images

```
/backup-poc.sh -u mysqlbackup -p new-password restore
```

### clean-up all old backups leaving only the current backup

```
/backup-poc.sh clean
```

---

## example scenario (assuming mysqlbackup user is initialised)

```
# create a new 'current' backup and take some deltas
/backup-poc.sh -u mysqlbackup -p new-password create_backup
/backup-poc.sh -u mysqlbackup -p new-password create_delta
/backup-poc.sh -u mysqlbackup -p new-password create_delta
/backup-poc.sh -u mysqlbackup -p new-password create_delta

# create a new 'current' backup and take some deltas - the previous one is now inactive
/backup-poc.sh -u mysqlbackup -p new-password create_backup
/backup-poc.sh -u mysqlbackup -p new-password create_delta
/backup-poc.sh -u mysqlbackup -p new-password create_delta

# restore from 'current' backup (snapshot + deltas)
/backup-poc.sh -u mysqlbackup -p new-password restore

# clean 'non-current' local backup files
/backup-poc.sh clean
```

source /backup-poc.sh
/backup-poc.sh -u root -p root init
/backup-poc.sh -u mysqlbackup -p new-password create_backup
/backup-poc.sh -u mysqlbackup -p new-password create_delta

/backup-poc.sh -u mysqlbackup -p new-password restore




