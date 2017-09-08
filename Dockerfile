FROM oraclelinux:7.3

COPY dist/bin/weblogic-operator /

ENTRYPOINT ["/weblogic-operator"]
