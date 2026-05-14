# Engineering Metrics

This document summarizes the intended measurement metrics for SLA-oriented network analysis.

## Delay metrics

| Metric | Meaning |
|---|---|
| Average latency | baseline transport delay |
| p95 latency | service consistency |
| p99 latency | congestion visibility |
| Maximum latency | worst-case behavior |

## Jitter metrics

| Metric | Meaning |
|---|---|
| RMS jitter | timing stability |
| Peak jitter | transient disturbance visibility |
| Jitter histogram | queue dynamics |

## Reliability metrics

| Metric | Meaning |
|---|---|
| Packet loss | transport integrity |
| Packet reorder | congestion side effects |
| Burst loss | instability events |

## Timestamp credibility

The repository emphasizes documenting:

- timestamp origin;
- synchronization assumptions;
- timing uncertainty;
- datapath visibility.

## Reporting philosophy

Metrics should be:

- reproducible;
- machine-readable;
- exportable to reports and dashboards.
