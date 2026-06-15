# SLA Measurement Report Example

## Executive summary

| Field | Value |
|---|---|
| Scenario | baseline network test with two SFP modules in a serial path |
| Measurement window | sample exported interval from `results/sample-test-1/` |
| SLA status | PASS |
| Main observation | stable microsecond-scale delay with negligible loss and near line-rate throughput |

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
| Timestamp origin | sample exported summary from the repository result package |
| Clock source | not modeled in detail in this example package |
| Uncertainty | low for documentation review, still not a substitute for calibrated hardware validation |
| Main risk | the sample package demonstrates reporting structure more strongly than clock-origin rigor |

## Example acceptance thresholds

| Metric | Threshold | Measured | Status |
|---|---:|---:|---|
| RTT avg | <= 15.0 us | 12.4 us | PASS |
| RTT max | <= 20.0 us | 16.8 us | PASS |
| Jitter | <= 1.0 us | 0.4 us | PASS |
| Packet loss | <= 0.01% | 0.0001% | PASS |
| Throughput | >= 9.5 Gbps | 9.8 Gbps | PASS |

## Available artifacts

- `results/sample-test-1/summary.md`
- `results/sample-test-1/plots/sample-metrics.svg`
- `docs/assets/sla_dashboard.svg`

## Delay and directionality

The committed sample package reports:

- RTT min: 11.7 us
- RTT avg: 12.4 us
- RTT max: 16.8 us
- OWD forward: 6.1 us
- OWD reverse: 6.3 us

The forward and reverse directions are close enough to treat the path as balanced in this example.

## Jitter and loss interpretation

The sample run shows low timing variation and effectively zero loss:

- jitter: 0.4 us
- packet loss: 0.0001%
- throughput: 9.8 Gbps

That combination is consistent with a clean baseline path rather than congestion, queue buildup or bursty impairment.

## Root-cause hypothesis

Potential explanations:

- no strong fault signal is visible in the sample package;
- residual variation is consistent with a healthy baseline export;
- deeper root-cause work would need packet-level traces and timestamp-origin details.

## Engineering conclusion

The committed result package is suitable as a reviewer-facing example because it contains a concrete summary, a concrete plot artifact and a filled report. It still does not replace calibrated hardware validation, but it closes the gap between pure templates and real measurement evidence.
