FROM centurylink/ca-certs
MAINTAINER Sercan Degirmenci <sercan@otsimo.com>

ADD bin/otsimo-catalog-linux-amd64 /opt/otsimo-catalog/bin/otsimo-catalog

# enable verbose debug for now
CMD ["/opt/otsimo-catalog/bin/otsimo-catalog","--debug","--storage","mongodb"]
