FROM alpine:3.5

RUN apk add --update ca-certificates git && rm -rf /var/cache/apk/*

ADD catalog-linux-amd64 /opt/otsimo/catalog

EXPOSE 18857

CMD ["/opt/otsimo/catalog"]
