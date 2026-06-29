# SLA Measurement Report Example

## Executive summary

| Field | Value |
|---|---|
| Scenario | baseline network test with two SFP modules in a serial path |
| Measurement window | sample exported interval from `results/sample-test-1/` |
| SLA status | PASS |
| Main observation | stable microsecond-scale delay with negligible loss and near line-rate throughput |
| Customer-facing conclusion | the sample path meets the example latency, jitter, loss and throughput acceptance gates |

## Customer decision

The measured sample package is suitable for a baseline acceptance report. No immediate network fault is visible in the example data. The recommended action is to preserve this run as the reference baseline and compare later field measurements against the same threshold table.

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

| Metric | Threshold | Measured | Margin | Status |
|---|---:|---:|---:|---|
| RTT avg | <= 15.0 us | 12.4 us | 2.6 us | PASS |
| RTT max | <= 20.0 us | 16.8 us | 3.2 us | PASS |
| Jitter | <= 1.0 us | 0.4 us | 0.6 us | PASS |
| Packet loss | <= 0.01% | 0.0001% | 0.0099% | PASS |
| Throughput | >= 9.5 Gbps | 9.8 Gbps | 0.3 Gbps | PASS |

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

## Limitations

This is a compact repository example, not a calibrated production acceptance test. A production report should additionally include:

- exact hardware model and firmware version;
- timestamp source and clock synchronization method;
- packet rate and payload size;
- test duration;
- raw or reduced trace checksum;
- operator and environment notes;
- uncertainty budget.

## Recommended follow-up

1. Run the synthetic SLA generator and confirm the report pipeline from a clean clone.
2. Replace sample metrics with a real trace package when hardware timestamp data is available.
3. Preserve this report structure for customer-facing acceptance documents.
4. Add a before/after comparison when the report is used for fault localization.

## Engineering conclusion

The committed result package is suitable as a reviewer-facing and customer-facing example because it contains a concrete summary, concrete plot artifacts, threshold margins and a filled report. It still does not replace calibrated hardware validation, but it closes the gap between pure templates and real measurement evidence.
