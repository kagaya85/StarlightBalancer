FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src/services/push

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates  \
    netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/services/push/bin /app

WORKDIR /app

EXPOSE 8083
EXPOSE 9093
VOLUME /data/conf

CMD ["./push", "-conf", "/data/conf"]
