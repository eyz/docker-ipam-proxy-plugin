FROM alpine

RUN apk update && \
    mkdir -p /run/docker/plugins

COPY docker-ipam-proxy-plugin docker-ipam-proxy-plugin
