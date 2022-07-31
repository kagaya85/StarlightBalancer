FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src/services/transcode

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates  \
    netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/services/transcode/bin /app

WORKDIR /app

EXPOSE 8085
EXPOSE 9095
VOLUME /data/conf

CMD ["./transcode", "-conf", "/data/conf"]
