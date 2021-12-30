# prometheus-vmware-exporter
Collect metrics ESXi Host

## Build
```bash 
docker build -t prometheus-vmware-exporter .
```

## Run
配置文件config.yaml
```
# tls_server_config:
#   cert_file: exporter.crt
#   key_file: exporter.key
# http_server_config:
#   # Enable HTTP/2 support. Note that HTTP/2 is only supported with TLS.
#   # This can not be changed on the fly.
#   [ http2: <bool> | default = true ]
basic_auth_users:
  # 当前设置的用户名为 prometheus ， 可以设置多个
  esxi: $2y$12$zwhHgEA9xxxxxxxxxxxxxxxxxxxxxxxxxx
```


```bash
docker run -tid \
  --restart=always \
  --name=prometheus-vmware-exporter \
  -e ESX_HOST=esx.domain.local \
  -e ESX_USERNAME=user \
  -e ESX_PASSWORD=password \
  -e ESX_LOG=debug \
  -e CONFIG=/etc/vmware/config/config.yaml \
  -v config.yaml:/etc/vmware/config/config.yaml \
  prometheus-vmware-exporter 
```

docker-compose.yaml

```yaml
version: "3"

networks:
  esxi:

services:
  exsi10:
    image: prometheus-vmware-exporter:0.2.1
    container_name: exsi10
    hostname: exsi10
    restart: always
    # privileged: true
    # user: root
    ports:
      - "9810:9879"
    volumes:
      - /etc/localtime:/etc/localtime:ro
      - /yaml/vmware/config/:/etc/vmware/config
    command: --config=/etc/vmware/config/config.yaml
    environment:
      - ESX_HOST=host
      - ESX_USERNAME=user
      - ESX_PASSWORD=password
      - ESX_LOG=info
    networks:
      - esxi
```