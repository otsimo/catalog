FROM alpine:3.4
MAINTAINER Sercan Degirmenci <sercan@otsimo.com>

RUN apk add --update ca-certificates git && rm -rf /var/cache/apk/*

ADD catalog-linux-amd64 /opt/otsimo/catalog

EXPOSE 18857

CMD ["/opt/otsimo/catalog"]
