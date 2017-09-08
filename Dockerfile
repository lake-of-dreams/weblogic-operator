FROM oraclelinux:7.3

COPY dist/bin/mysql-operator /

ENTRYPOINT ["/mysql-operator"]
