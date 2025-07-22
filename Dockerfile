FROM golang:latest as builder
WORKDIR /app
COPY . /app/
ENV GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 

RUN go mod tidy && go build -ldflags="-w -s"


FROM debian:stable-slim

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list  \
    && apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates  \
    netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /app/prometheus-vmware-exporter /usr/bin/prometheus-vmware-exporter
EXPOSE 9879
ENTRYPOINT ["prometheus-vmware-exporter"]
