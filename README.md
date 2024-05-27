# Teeworlds Exporter for Prometheus

![golangci-lint](https://github.com/theobori/teeworlds-prometheus-exporter/actions/workflows/lint.yml/badge.svg) ![build](https://github.com/theobori/teeworlds-prometheus-exporter/actions/workflows/build.yml/badge.svg)

This is a server that scrapes Teeworlds master servers and econ servers, get stats and exports them via HTTP for Prometheus consumption.

## 📖 Build and run

For the build, you only need the following requirements:

- [Go](https://golang.org/doc/install) 1.22.3


Next to the Go application you could need the following requirements:
- [Teeworlds server](https://www.teeworlds.com/?page=downloads&id=14786) 0.7
  - With a econ server

Now you can build and run the Go application, check the `-h` or `--help` flag if needed.


## 🤝 Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).

## ⚙️ YAML configuration

The YAML configuration is supposed to have the following format.

```yaml
servers:
  econ:
    - host: localhost
      port: 7000
      password: hello_world

  master:
    - protocol: http
      url: "https://master1.ddnet.tw/ddnet/15/servers.json"
      refresh_cooldown: 10

    - protocol: udp
      host: "master1.teeworlds.com"
      port: 8283
      refresh_cooldown: 15
    
    - protocol: udp
      host: "master2.teeworlds.com"
      port: 8283
      refresh_cooldown: 15
    
    - protocol: udp
      host: "master3.teeworlds.com"
      port: 8283
      refresh_cooldown: 15

    - protocol: udp
      host: "master4.teeworlds.com"
      port: 8283
      refresh_cooldown: 15

```