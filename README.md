A Cisco IronPort exporter for Prometheus. Production ready.

# Getting started

Download and install:

```
go get github.com/1player/ironport-prometheus-exporter
```

Run the exporter:

```
ironport-prometheus-exporter -ironport.host 192.168.0.1 -ironport.basic-auth user:pass -web.listen-address 0.0.0.0:9101
```

Help on flags:
```
ironport-prometheus-exporter --help
```

# Gotchas

Since I have no idea what each IronPort statistic actually measures, the HELP labels on the metrics are unset.
