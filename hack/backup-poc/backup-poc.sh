#!/bin/bash

MYSQL_CONFIG=/etc/my.cnf
MYSQL_DATA_DIR=/var/lib/mysql
SOCKET=/var/lib/mysql/mysql.sock


ADMIN_USER=root
ADMIN_USER_PASSWORD=root

BACKUP_USER=mysqlbackup
BACKUP_USER_PASSWORD=new-password

BACKUP_BASE_DIR=/backups

SS_IMG_NAME=snapshot-image.mbi
SS_OUT_NAME=create-snapshot.out
SSR_OUT_NAME=restore-snapshot.out
INC_IMG_NAME=ind-delta-image.bi
ID_OUT_NAME=create-inc-delta.out

# -------------------------------------------------------------------------------------------------
# general functions

function log() { 
    echo "$@"; 
}

function log_err() { 
    echo "$@" 1>&2; 
}

# -------------------------------------------------------------------------------------------------
# mysql init functions

function create_mysqlbackup_user() {
    local admin_user=${1:-${ADMIN_NAME}}
    local admin_user_password=${2:-${ADMIN_USER_PASSWORD}}
    mysql -u${admin_user} -p${admin_user_password} < backup-poc-create-user.sh
}

# -------------------------------------------------------------------------------------------------
# filesystem functions

function create_filesystem_paths() {
    local paths=("$@")
    for path in "${paths[@]}"; do
        log "created: ${path}"
        mkdir -p ${path}
    done
}

# NB : does not deal with spaces in file names.
function list-directory-files() {
    local target_dir=$1
    local files=`ls -l ${target_dir} | grep '^-' | awk -F " " '{print $NF}'`
    echo ${files}
}

# NB : does not deal with spaces in file names.
function list-directory-file-paths() {
    local target_dir=$1
    local file_paths=`ls -l ${target_dir} | grep '^-' | awk -F " " -v path="${target_dir}/" '{a=path$NF; print a}'`
    echo ${file_paths}
}

# NB : does not deal with spaces in file names.
function list-directory-dirs() {
    local target_dir=$1
    local files=`ls -l ${target_dir} | grep '^d' | awk -F " " '{print $NF}'`
    echo ${files}
}

# NB : does not deal with spaces in file names.
function list-directory-dir-paths() {
    local target_dir=$1
    local file_paths=`ls -l ${target_dir} | grep '^d' | awk -F " " -v path="${target_dir}/" '{a=path$NF; print a}'`
    echo ${file_paths}
}

# NB: assumes the specified directory contains a single file
# NB : does not deal with spaces in file names.
function get-single-file() {
    local target_dir=$1
    local single_file=`list-directory-files ${target_dir}`
    # TODO: do check - error if no single file
    echo ${single_file}
}

# -------------------------------------------------------------------------------------------------
# local backup filesystem functions

# clean-up all local backup files except for the 'current' backup.
function cleanup_old_backups() {
    local backup_paths=(`list-directory-dir-paths ${BACKUP_BASE_DIR}`)
    unset "backup_paths[${#backup_paths[@]}-1]"
    for backup_path in "${backup_paths[@]}"; do
        echo "removing non-current backup: ${backup_path}"
        delete_backup ${backup_path}
    done
}

# default deletes the specified backup.
function delete_backup() {
    local backup_dir=$1
    rm -Rf ${backup_dir}
}

# -------------------------------------------------------------------------------------------------
# backup meta functions

# generate a time-based lexographically ordered 'id' to identify backups, snapshots, deltas, and,
# other artifacts.
function generate_timestamp_id() {
    local id=`date +%Y-%m-%d-%H%M%S-%N`
    echo ${id}
}

# NB : assumes snapshots directory exists.
# NB : does not deal with spaces in directory names.
function list_backup_ids() {
    local backup_ids=`ls -l ${BACKUP_BASE_DIR} | grep '^d' | awk -F " " '{print $NF}'`
    echo "${backup_ids}"
}

# NB : assumes snapshots created with option '--with-timestamp' or are otherwise lexigraphically named.
function get_current_backup_id() {
    local last=`list_backup_ids | sort | tail -1`
    echo ${last}
}

# NB : assumes snapshots directory exists.
# NB : does not deal with spaces in directory names.
function list_backup_paths() {
    local backup_paths=`ls -l ${BACKUP_BASE_DIR} | grep '^d' | awk -F " " -v path="${BACKUP_BASE_DIR}/" '{a=path$NF; print a}'`
    echo "${backup_paths}"
}

# NB : assumes snapshots created with option '--with-timestamp' or are otherwise lexigraphically named.
function get_current_backup_path() {
    local current_backup_path=`list_backup_paths | sort | tail -1`
    echo ${current_backup_path}
}

function get_current_backup_snapshot_path() {
    local current_backup_path=`get_current_backup_path`
    local snapshot_dir=`ls -l ${current_backup_path} | grep '^d' | grep snapshot | awk -F " " '{print $NF}'`
    echo ${current_backup_path}/${snapshot_dir}
}

function echo_current_backup_snapshot_info() {
    local current_backup_snapshot_path=`get_current_backup_snapshot_path`
    cat ${current_backup_snapshot_path}/backup/meta/backup_variables.txt
}

function get_current_backup_snapshot_value() {
    local key=$1
    local value=`echo_current_backup_snapshot_info | grep ${key} | cut -d'=' -f2`
    echo ${value}
}

function get_current_backup_snapshot_meta_value() {
    local value=`get_current_backup_snapshot_value start_lsn`
    echo ${value}
}

function list_current_backup_inc_delta_paths() {
    local current_backup_path=`get_current_backup_path`
    local backup_inc_delta_paths=`ls -l ${current_backup_path} | grep '^d' | grep inc-delta | awk -F " " -v path="${current_backup_path}/" '{a=path$NF; print a}'`
    echo ${backup_inc_delta_paths}
}

# -------------------------------------------------------------------------------------------------
# mysqlbackup output validation functions

function output_assert_exists() {
    local output=$1
    if [[ -z ${output} ]]; then
        log_err "no output filepath was specified."
        exit 1
    fi
    if [[ ! -f ${output} ]]; then
        echo "no output file exists was: ${output}"
        exit 1
    fi
}

function output() {
    local output=$1
    output_assert_exists ${output}
    local output=`cat ${output}`
    echo "${output}"
}

function output_warnings() {
    local output=$1
    output_assert_exists ${output}
    local warnings=`cat ${output} | grep 'MAIN WARNING:'`
    echo "${warnings}"
}

function output_errors() {
    local output=$1
    echo "output errors"
    output_assert_exists ${output}
    local errors=`cat ${snapshot_output} | grep 'MAIN ERROR:'`
    echo outputting errors
    echo "${errors}"
}

function output_failure() {
    local output=$1
    output_assert_exists ${output}
    local failed=`cat ${output} | grep '^mysqlbackup failed with errors!'`
    echo "${failed}"
}

function output_success() {
    local output=$1
    output_assert_exists ${output}
    local success=`cat ${output} | grep '^mysqlbackup completed OK!'`
    echo ${success}
}

function output_start_lsn() {
    local output=$1
    output_assert_exists ${output}
    local start_lsn=`cat ${output} | grep 'Start LSN' | cut -d':' -f2 | tr -d '[:space:]'`
    echo ${start_lsn}
}

function output_end_lsn() {
    local output=$1
    output_assert_exists ${output}
    local end_lsn=`cat ${output} | grep 'End LSN' | cut -d':' -f2 | tr -d '[:space:]'`
    echo ${end_lsn}
}

function output_assert_success() {
    local output=$1
    output_assert_exists ${output}
    local failed=`output_failure ${output}`
    if [[ ${failed} ]]; then
        log_err "backup operation failed failed: output failure asserted."
        echo "loggin errors"
        echo "${output}"
        log_err `output_errors ${output}`
        exit 1
    fi
    local success=`output_success ${output}`
    if [[ -z ${success} ]]; then
        log_err "backup operation failed: output success not asserted."
        # log_err `output_errors ${output}`
        exit 1
    fi
}

# -------------------------------------------------------------------------------------------------
# snapshot output validation functions

function echo_mysql_backup_history() {
    local backup_user=${1:-${BACKUP_USER}}
    local backup_user_password=${2:-${BACKUP_USER_PASSWORD}}
        mysql \
        -u${backup_user} \
        -p${backup_user_password} \
        -e 'select * from mysql.backup_history\G'
}

# -------------------------------------------------------------------------------------------------
# backup functions

function init_new_backup() {
    if [ $# -ne 6 ]; then
        log_err "incorrect number of arguments passes to 'init_new_snapshot': $@"
        exit 1
    fi

    local backup_id=$1
    local snapshot_name=${backup_id}-snapshot
    local _snapshot_dir=${BACKUP_BASE_DIR}/${backup_id}/${snapshot_name}
    local _snapshot_backup_dir=${_snapshot_dir}/backup
    local _snapshot_image_dir=${_snapshot_dir}/backup-image
    local _snapshot_img=${_snapshot_image_dir}/${snapshot_name}_${SS_IMG_NAME}
    local _snapshot_ouput=${_snapshot_dir}/${SS_OUT_NAME}

    create_filesystem_paths ${_snapshot_dir} ${_snapshot_backup_dir} ${_snapshot_image_dir}

    local __return_var_snapshot_dir=$2
    eval $__return_var_snapshot_dir="'$_snapshot_dir'"
    local __return_var_snapshot_backup_dir=$3
    eval $__return_var_snapshot_backup_dir="'$_snapshot_backup_dir'"
    local __return_var_snapshot_image_dir=$4
    eval $__return_var_snapshot_image_dir="'$_snapshot_image_dir'"
    local __return_var_snapshot_img=$5
    eval $__return_var_snapshot_img="'$_snapshot_img'"
    local __return_var_snapshot_ouput=$6
    eval $__return_var_snapshot_ouput="'$_snapshot_ouput'"
}

function create_new_backup() {
    local backup_id=`generate_timestamp_id`
    log "create_new_backup: ${backup_id}"
    local backup_user=${1:-${BACKUP_USER}}
    local backup_user_password=${2:-${BACKUP_USER_PASSWORD}}

    init_new_backup ${backup_id} \
        snapshot_dir snapshot_backup_dir snapshot_image_dir snapshot_img snapshot_ouput

    mysqlbackup \
        --defaults-file=${MYSQL_CONFIG} \
        --user=${backup_user} \
        --password=${backup_user_password} \
        --socket=${SOCKET} \
        --backup-dir=${snapshot_backup_dir} \
        --backup-image=${snapshot_img} \
        backup-to-image 2> ${snapshot_ouput}
    
    log "finished backup"

    log "mysqlbackup output: ${snapshot_ouput}"
    output_assert_success ${snapshot_ouput}
    output_warnings ${snapshot_ouput}

    local start_lsn=`output_start_lsn ${snapshot_ouput}`
    local end_lsn=`output_end_lsn ${snapshot_ouput}`
    echo ${start_lsn}:${end_lsn}
}

function create_next_inc_delta() {
    log "create_next_inc_delta: `get_current_backup_id`"
    local backup_user=${1:-${BACKUP_USER}}
    local backup_user_password=${2:-${BACKUP_USER_PASSWORD}}

    local delta_id=`generate_timestamp_id`
    local backup_path=`get_current_backup_path`
    local delta_dir=${backup_path}/${delta_id}-inc-delta
    local delta_backup_dir=${delta_dir}/backup
    local delta_img_dir=${delta_dir}/backup-image
    local delta_img_name=${delta_id}-${INC_IMG_NAME}
    local delta_img=${delta_img_dir}/${delta_img_name}
    local delta_ouput=${delta_dir}/${ID_OUT_NAME}

    create_filesystem_paths ${delta_dir} ${delta_backup_dir} ${delta_img_dir} 

    mysqlbackup \
        --defaults-file=${MYSQL_CONFIG} \
        --user=${backup_user} \
        --password=${backup_user_password} \
        --socket=${SOCKET} \
        --incremental \
        --incremental-base=history:last_backup \
        --backup-dir=${delta_backup_dir} \
        --backup-image=${delta_img} \
        backup-to-image 2> ${delta_ouput}

    log "mysqlbackup output: ${delta_ouput}"
    output_assert_success ${delta_ouput}
    output_warnings ${delta_ouput}

    local start_lsn=`output_start_lsn ${delta_ouput}`
    local end_lsn=`output_end_lsn ${delta_ouput}`
    echo ${start_lsn}:${end_lsn}
}

# -------------------------------------------------------------------------------------------------
# image validation functions

function validate_image() {
    log "validate_image"
    
    local backup_img=$1
    local tmp_dir="/tmp/validate"
    local validation_ouput=${tmp_dir}/`generate_timestamp_id`-validation.out
    
    create_filesystem_paths ${tmp_dir}
    
    mysqlbackup \
        --backup-image=${validation_ouput} \
		validate 2> /tmp/
    snapshot_assert_success ${validation_ouput}

    log "mysqlbackup output: ${validation_ouput}"
    output_assert_success ${validation_ouput}
    output_warnings ${validation_ouput}
}

# -------------------------------------------------------------------------------------------------
# restore functions

function restore_current_backup() {
    log "restore_current_backup: `get_current_backup_id`"
    local backup_user=${1:-${BACKUP_USER}}
    local backup_user_password=${2:-${BACKUP_USER_PASSWORD}}

    local backup_snapshot_path=`get_current_backup_snapshot_path`
    log "restoring snapshot: ${backup_snapshot_path}"
    restore_backup_snapshot ${backup_user} ${backup_user_password} ${backup_snapshot_path}

    local inc_delta_paths=(`list_current_backup_inc_delta_paths`)
    for inc_delta_path in "${inc_delta_paths[@]}"; do
        log "restoring inc_delta: ${inc_delta_path}"
        restore_backup_next_inc_delta ${backup_user} ${backup_user_password} ${inc_delta_path}
    done

    log "restored backup: `get_current_backup_id`"
}

function restore_backup_snapshot() {
    log "restore_backup_snapshot: `get_current_backup_id`"
    local backup_user=${1:-${BACKUP_USER}}
    local backup_user_password=${2:-${BACKUP_USER_PASSWORD}}
    local backup_snapshot_dir=$3

    local backup_snapshot_img_dir=${backup_snapshot_dir}/backup-image
    local backup_snapshot_img=${backup_snapshot_img_dir}/`get-single-file ${backup_snapshot_img_dir}`
    local restore_snapshot_dir=${backup_snapshot_dir}/restore
    local restore_ouput=${backup_snapshot_dir}/${SSR_OUT_NAME}

    create_filesystem_paths ${restore_snapshot_dir}

    # NB: remove force and test properly? mysql might be in sidecar or local to container...
    mysqlbackup \
        --defaults-file=${MYSQL_CONFIG} \
        --user=${backup_user} \
        --password=${backup_user_password} \
        --socket=${SOCKET} \
        --datadir=${MYSQL_DATA_DIR} \
        --backup-dir=${restore_snapshot_dir} \
        --backup-image=${backup_snapshot_img} \
        copy-back-and-apply-log  2> ${restore_ouput}

    log "mysqlbackup output: ${restore_ouput}"
    output_assert_success ${restore_ouput}
    output_warnings ${restore_ouput}
}

function restore_backup_next_inc_delta() {
    log "restore_backup_next_inc_delta: `get_current_backup_id`"
    local backup_user=${1:-${BACKUP_USER}}
    local backup_user_password=${2:-${BACKUP_USER_PASSWORD}}
    local backup_inc_delta_dir=$3
    
    local backup_inc_delta_img_dir=${backup_inc_delta_dir}/backup-image
    local backup_inc_delta_img=${backup_inc_delta_img_dir}/`get-single-file ${backup_inc_delta_img_dir}`
    local restore_inc_delta_dir=${backup_inc_delta_dir}/restore
    local restore_ouput=${backup_inc_delta_dir}/${SSR_OUT_NAME}

    create_filesystem_paths ${restore_snapshot_dir}

    mysqlbackup \
        --defaults-file=${MYSQL_CONFIG} \
        --user=${backup_user} \
        --password=${backup_user_password} \
        --socket=${SOCKET} \
        --incremental \
        --datadir=${MYSQL_DATA_DIR} \
        --incremental-backup-dir=${restore_inc_delta_dir} \
        --backup-image=${backup_inc_delta_img} \
        --force \
        copy-back-and-apply-log  2> ${restore_ouput}

    log "mysqlbackup output: ${restore_ouput}"
    output_assert_success ${restore_ouput}
    output_warnings ${restore_ouput}
}

# -------------------------------------------------------------------------------------------------
# archive functions

function download_mysql_snapshot_backup() {
    # curl -s \
    #     -u ‘timothy.langford@oracle.com:XXX’ \
    #     https://swiftobjectstorage.us-phoenix-1.oraclecloud.com/v1/bristol_dev/backups/backup.img \
    #     -o backup.img
    exit 2
}

function archive_mysql_snapshot_backup() {
    # curl \
    #     -s \
    #     -X PUT \
    #     -T backup.img \
    #     -u ‘timothy.langford@oracle.com:XXX’ \
    #     https://swiftobjectstorage.us-phoenix-1.oraclecloud.com/v1/bristol_dev/backups/backup.img
    exit 2
}

# -------------------------------------------------------------------------------------------------
# main functions

function display_usage() {
    echo 'usage : backup-poc -u ${mysql_user} -p ${mysql-passwd} ${command}'
    echo 'where - ${command} : create_snapshot | create_differential_delta | create_incremental_delta'
}

num_params=$#
while [[ $# -gt 1 ]]; do
option_key="$1"
option_value="$2"
    case ${option_key} in
        -u|--user)
        mysql_user=${option_value}
        shift 2
        ;;
        -p|--password)
        mysql_passwd=${option_value}
        shift 2
        ;;
        *)
        echo "error: unknown option: ${option_key} ${option_value}"
        break
        ;;
    esac
done

# the command is the last argument
backup_cmd=$1
case ${backup_cmd} in
    init)
    echo "command: create_user"
    create_mysqlbackup_user ${mysql_user} ${mysql_passwd}
    ;;
    create_backup)
    echo "command: create_backup"
    create_new_backup ${mysql_user} ${mysql_passwd}
    ;;
    create_delta)
    echo "command: create_delta"
    create_next_inc_delta ${mysql_user} ${mysql_passwd}
    ;;
    restore)
    echo "command: restore"
    restore_current_backup ${mysql_user} ${mysql_passwd}
    ;;
    clean)
    echo "command: clean"
    cleanup_old_backups
    ;;
    "")
    if [[ ${num_params} -gt 1 ]]; then
        echo "error: no command was specified"
    fi
    ;;
    *)
    echo "error: unknown command: ${backup_cmd}"
    ;;
esac
