# Demo report

## Scenario

Hardware-assisted SLA validation demo with synthetic benchmark data.

## Input

```bash
python tools/generate_demo_benchmark.py
```

## Output artifacts

```text
results/demo-benchmark/metrics.csv
docs/assets/generated_benchmark.svg
```

## Acceptance summary

| Metric | Target | Demo result | Status |
|---|---:|---:|---|
| Packet loss | <= 0.10 % | 0.02 % | PASS |
| One-way delay | <= 1.000 ms | 0.384 ms | PASS |
| Jitter | <= 100 us | 42 us | PASS |
| Throughput | >= 900 Mbit/s | 941 Mbit/s | PASS |

## Interpretation

The demo result represents a stable service profile:

- low packet loss;
- sub-ms one-way delay;
- bounded jitter;
- throughput close to line rate.

## Root-cause note

No critical SLA violation is detected in the demo profile.

A real production run would additionally correlate jitter, delay, loss and throughput with hardware timestamp origin, PTP status and topology metadata.
