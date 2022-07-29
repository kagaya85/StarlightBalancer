FROM golang:1.18 AS builder

COPY . /src
WORKDIR /src/services/timeline

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates  \
    netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/services/timeline/bin /app

WORKDIR /app

EXPOSE 8087
EXPOSE 9097
VOLUME /data/conf

CMD ["./timeline", "-conf", "/data/conf"]
