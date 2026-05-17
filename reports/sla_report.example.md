# SLA Measurement Report Example

## Executive summary

| Field | Value |
|---|---|
| Scenario | software vs datapath timing comparison |
| Measurement window | example synthetic interval |
| SLA status | TBD |
| Main observation | timestamp origin strongly affects confidence |

## Measurement topology

```text
probe generator
-> network path
-> timestamp point
-> metric extraction
-> SLA report
```

## Timing model

| Item | Description |
|---|---|
| Timestamp origin | software or datapath |
| Clock source | documented per scenario |
| Uncertainty | explicitly listed in the report |
| Main risk | OS noise or clock drift, depending on method |

## SLA thresholds

| Metric | Threshold | Measured | Status |
|---|---:|---:|---|
| p95 latency | TBD | TBD | TBD |
| p99 latency | TBD | TBD | TBD |
| RMS jitter | TBD | TBD | TBD |
| Packet loss | TBD | TBD | TBD |

## Latency distribution

Expected artifact:

```text
reports/latency_distribution.svg
```

## Jitter analysis

Expected artifact:

```text
reports/jitter_histogram.svg
```

## Packet-loss timeline

Expected artifact:

```text
reports/packet_loss_timeline.svg
```

## Root-cause hypothesis

Potential explanations:

- queue buildup;
- unstable host timing;
- packet burst effects;
- clock uncertainty;
- link-level issues.

## Engineering conclusion

The report should state whether the SLA conclusion is valid, uncertain or blocked by missing timing assumptions.
