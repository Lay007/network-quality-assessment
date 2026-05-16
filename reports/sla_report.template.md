# SLA Measurement Report

## Executive summary

| Field | Value |
|---|---|
| Scenario | TBD |
| Measurement window | TBD |
| SLA status | PASS / FAIL / INCONCLUSIVE |
| Main risk | TBD |

## Measurement topology

```text
traffic source
-> network path
-> timestamp point
-> metrics extraction
-> report
```

## Timing model

Document:

- timestamp origin;
- clock source;
- synchronization assumption;
- timestamp uncertainty;
- packet serialization assumptions.

## SLA thresholds

| Metric | Threshold | Measured | Status |
|---|---:|---:|---|
| p95 latency | TBD | TBD | TBD |
| p99 latency | TBD | TBD | TBD |
| RMS jitter | TBD | TBD | TBD |
| Packet loss | TBD | TBD | TBD |

## Latency distribution

Include histogram or percentile plot.

## Jitter analysis

Explain whether jitter is stable, bursty or correlated with packet-loss events.

## Packet loss timeline

Document:

- isolated losses;
- burst losses;
- duration of loss events;
- possible correlation with delay spikes.

## Root-cause hypothesis

List the most likely causes:

- queue buildup;
- congestion;
- clock drift;
- software timestamp noise;
- physical link issue.

## Engineering conclusion

Summarize what the measurement proves, what remains uncertain and what should be tested next.
