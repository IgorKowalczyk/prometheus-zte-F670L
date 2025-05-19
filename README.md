# Prometheus exporter for ZTE F670L ONT

Exporter for ZTE F670L ONT, based on the work of [lucathehacker](https://github.com/LucaTheHacker/PrometheusF6005).

## Required environment variables

| Name             | Description                  | Default Value      |
| ---------------- | ---------------------------- | ------------------ |
| `ENDPOINT`       | HTTP address to the ONT      | http://192.168.1.1 |
| `ONT_USERNAME`   | Username for the ONT         | user               |
| `ONT_PASSWORD`   | Password for the ONT         | user               |
| `ONT_SLEEP_QUIT` | Time to wait before quitting | 60s                |

## Usage

Example docker-compose section:

```yaml
f6005_exporter:
 restart: always
 image: ghcr.io/igorkowalczyk/prometheus-zte-F670L
 environment:
  - ENDPOINT=http://192.168.1.1
  - ONT_USERNAME=user
  - ONT_PASSWORD=user
 expose:
  - 80
```
