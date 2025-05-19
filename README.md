# Prometheus Exporter for ZTE F670L ONT

A Prometheus exporter for ZTE F670L ONT, enabling you to monitor device metrics such as CPU, memory, network statistics, and connected clients.

<!-- prettier-ignore-start -->
> [!IMPORTANT]
> **This exporter is only tested with the `V9.0.11P1N2` firmware.**
> Other firmware versions may work, but are not guaranteed. If you have a different version and it works, please let me know or create a pull request.
<!-- prettier-ignore-end -->

## ⚙️ Configuration

Set the following environment variables (defaults shown):

| Name             | Description                          | Default Value      |
| ---------------- | ------------------------------------ | ------------------ |
| `ENDPOINT`       | HTTP address of the ONT              | http://192.168.1.1 |
| `ONT_USERNAME`   | Username for the ONT                 | `user`             |
| `ONT_PASSWORD`   | Password for the ONT                 | `user`             |
| `ONT_SLEEP_QUIT` | Seconds to wait before exit on error | `60`               |

You can set these in your environment, `.env` file, or directly in the `docker-compose.yml` file.

---

## 📦 `docker-compose.yml` file.

```yaml
services:
 f670l_exporter:
  image: ghcr.io/igorkowalczyk/prometheus-zte-f670l:latest
  restart: always
  environment:
   - ENDPOINT=http://192.168.1.1
   - ONT_USERNAME=user
   - ONT_PASSWORD=user
  ports:
   - 3000:3000
```

The exporter will be available at `http://localhost:3000/metrics`.

## 🔥 Prometheus configuration

Add a scrape job to your `prometheus.yml` configuration file:

```yaml
scrape_configs:
 - job_name: "zte-f670l"
   static_configs:
    - targets: ["localhost:3000"]
```

<!-- prettier-ignore-start -->
> [!NOTE]
> Make sure to replace `localhost` with the actual IP address of the machine running the exporter if it's not on the same machine as Prometheus.
<!-- prettier-ignore-end -->

## 📊 Grafana

A Grafana dashboard is available for visualizing the metrics. https://grafana.com/grafana/dashboards/23453

[![Grafana Example](https://github.com/user-attachments/assets/c20d1fc1-4def-4943-bbc0-15d224d3d970)](https://grafana.com/grafana/dashboards/23453)

## 💻 Build from source

```sh
git clone https://github.com/igorkowalczyk/prometheus-zte-F670L.git
cd prometheus-zte-F670L
go build -o prometheus_exporter .
./prometheus_exporter
```

## ⁉️ Troubleshooting

- **Logout from web interface**: ZTE Designed the web interface in a way that only one session can be active at a time. If you log in to the web interface, the exporter will be logged out, causing it to fail. Exporter will exit after 60 seconds by default and _should_ be restarted automatically by Docker. You can also set `ONT_SLEEP_QUIT` to a lower value to make it exit faster.
- **Login errors**: Check your credentials and ensure that the `ONT_USERNAME` and `ONT_PASSWORD` environment variables are set correctly. The default credentials are `user:user`. Also ensure that the ONT is reachable from the machine running the exporter.
- **Metrics not showing**: Ensure that the exporter is running and accessible. Check the logs for any errors or warnings.

If you have any other issues, please check the [issues](https://github.com/igorkowalczyk/prometheus-zte-F670L/issues) or create a new one.

## 📥 Pull Requests

When submitting a pull request:

- Clone the repository (`git clone https://github.com/igorkowalczyk/prometheus-zte-F670L.git`).
- Create a branch off of `master` and give it a meaningful name (e.g. `my-awesome-new-feature`).
- Open a [pull request](https://github.com/igorkowalczyk/majo.exe/pulls) on [GitHub](https://github.com) and describe the feature or fix.

## 💝 Credits

- [lucathehacker/PrometheusF6005](https://github.com/LucaTheHacker/PrometheusF6005) (inspiration and the base for this project)
