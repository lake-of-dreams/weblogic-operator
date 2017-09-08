FROM oraclelinux:7.3

# dev mode tools
# ENV http_proxy http://www-proxy.us.oracle.com:80
# ENV https_proxy http://www-proxy.us.oracle.com:80
# RUN yum install -y vim

ADD rpm/meb-4.1.0-rhel5.x86_64.rpm /
RUN rpm -i meb-4.1.0-rhel5.x86_64.rpm
RUN mv /opt/mysql/meb-4.1/bin/mysqlbackup /sbin/mysqlbackup
RUN rm -R /opt/mysql
RUN rm meb-4.1.0-rhel5.x86_64.rpm

