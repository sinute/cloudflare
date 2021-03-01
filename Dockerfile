FROM alpine:latest

LABEL maintainer="sinute@outlook.com"

ADD cloudflare /bin
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.33-r0/glibc-2.33-r0.apk && \
    apk add glibc-2.33-r0.apk && \
    chmod 755 /bin/cloudflare
ENTRYPOINT ["/bin/cloudflare"]
