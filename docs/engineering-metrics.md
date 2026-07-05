# Engineering Metrics

This document summarizes the intended measurement metrics for SLA-oriented network analysis.

## Delay metrics

| Metric | Meaning | Acceptance note |
|---|---|---|
| Average latency | baseline transport delay | useful for trend, not enough for SLA proof alone |
| p95 latency | service consistency | good first percentile for user-visible degradation |
| p99 latency | congestion visibility | important for tail-latency analysis |
| Maximum latency | worst-case behavior | must be interpreted with trace duration and outlier notes |

## Jitter metrics

| Metric | Meaning | Acceptance note |
|---|---|---|
| RMS jitter | timing stability | useful when delay distribution is roughly stable |
| Peak jitter | transient disturbance visibility | should be tied to timestamp source and outlier context |
| Jitter histogram | queue dynamics | preferred for explaining whether jitter is random, bursty or multimodal |

## Reliability metrics

| Metric | Meaning | Acceptance note |
|---|---|---|
| Packet loss | transport integrity | should be derived from sequence accounting |
| Packet reorder | congestion side effects | should distinguish reorder from true loss |
| Burst loss | instability events | should include event length and time position |

## Timestamp credibility

The repository emphasizes documenting:

- timestamp origin;
- synchronization assumptions;
- timing uncertainty;
- datapath visibility.

A metric without timestamp provenance should be treated as diagnostic, not as strong SLA evidence.

## Minimum report table

Every promoted measurement report should include:

| Field | Why it matters |
|---|---|
| trace file | ties plots and metrics to input data |
| timestamp source | separates host artifacts from datapath behavior |
| clock assumption | explains whether one-way delay is credible |
| packet count | gives statistical context |
| loss count | supports reliability claims |
| p95 / p99 delay | shows tail behavior |
| jitter summary | shows timing stability |
| limitations | prevents overclaiming |

## Reporting philosophy

Metrics should be:

- reproducible;
- machine-readable;
- exportable to reports and dashboards;
- interpreted with limitations rather than presented as isolated numbers.
