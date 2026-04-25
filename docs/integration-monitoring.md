# Monitoring and integration

## Exported metrics

Typical metrics exposed:

- sla_delay_ms
- sla_jitter_us
- sla_loss_percent
- sla_throughput_mbps

## Prometheus integration

Example:

```text
sla_delay_ms 0.384
sla_jitter_us 42
sla_loss_percent 0.02
```

## Grafana dashboards

Suggested panels:

- delay over time
- jitter histogram
- SLA pass/fail status
- anomaly events timeline

## Alerting

Example rules:

```text
if sla_loss_percent > 0.1 → CRITICAL
if sla_jitter_us > 100 → WARNING
if sla_delay_ms > 1 → CRITICAL
```

## External integration

- syslog
- webhook
- REST API

## Engineering insight

Integration turns measurement into operational observability system.
