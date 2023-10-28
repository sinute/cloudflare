FROM golang:1.19.2-alpine3.16 as go-builder

WORKDIR /cloudflare

COPY . ./

RUN go mod verify
RUN go build -o bin/cloudflare main.go

FROM alpine:3.16

WORKDIR /cloudflare

LABEL maintainer="Sinute <sinute@outlook.com>"

COPY --from=go-builder /cloudflare/bin/cloudflare ./

RUN apk add --no-cache tzdata

ENTRYPOINT [ "/cloudflare/cloudflare" ]
