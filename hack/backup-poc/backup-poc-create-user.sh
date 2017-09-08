-- https://dev.mysql.com/doc/mysql-enterprise-backup/4.1/en/mysqlbackup.privileges.html/
--
create user 'mysqlbackup'@'localhost' identified by 'new-password';
grant RELOAD on *.* to 'mysqlbackup'@'localhost';
grant create, insert, drop, update on mysql.backup_progress to 'mysqlbackup'@'localhost';
grant create, insert, select, drop, update on mysql.backup_history to 'mysqlbackup'@'localhost';
grant replication client on *.* to 'mysqlbackup'@'localhost';
grant super on *.* to 'mysqlbackup'@'localhost';
grant process on *.* to 'mysqlbackup'@'localhost';
grant lock tables, select, create, drop, file on *.* to 'mysqlbackup'@'localhost';
grant create, insert, drop, update on mysql.backup_sbt_history to 'mysqlbackup'@'localhost';
-- grant alter on mysql.backup_history to 'mysqlbackup'@'localhost';
flush privileges;