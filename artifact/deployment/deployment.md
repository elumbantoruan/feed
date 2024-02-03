# Deployments artifact

## grafana

binary: /usr/share/grafana/bin/grafana  
config: /etc/grafana/grafana.ini  
systemd: /lib/systemd/system/grafana-server.service  

## tempo

binary: /usr/bin/tempo  
config: /etc/tempo/config.yml  
systemd: /etc/systemd/system/tempo.service

## prometheus

binary: /opt/prometheus/prometheus  
config: /opt/prometheus/prometheus.yml  
systemd: /etc/systemd/system/prometheus.service

## promtail

Deploy into K8S cluster

## loki
binary: /usr/bin/loki  
config: /etc/loki/lokiconfig.yaml  
systemd: /etc/systemd/system/loki.service

## minio

binary: /usr/local/bin/minio  
EnvironmentFile=-/etc/default/minio
systemd: /lib/systemd/system/minio.service