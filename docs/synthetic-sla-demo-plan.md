# Synthetic SLA Demo Plan

This document defines a hardware-free demo path for `network-quality-assessment`. The goal is to make the repository reproducible even when FPGA/SFP timestamping hardware is not available.

## Why a synthetic demo is needed

The repository presents a hardware timestamping methodology, but reviewers should be able to run a small experiment locally and see the complete reporting flow:

```text
synthetic packet trace -> SLA analysis -> summary CSV -> plots -> report
```

This does not replace real FPGA/SFP timestamping. It provides a public, deterministic baseline for documentation, CI and portfolio review.

## Proposed directory layout

```text
tools/
  generate_synthetic_sla_trace.py
  analyze_sla_trace.py
verification/reports/synthetic_sla_demo/
  synthetic_trace.csv
  sla_summary.csv
  one_way_delay_timeseries.png
  jitter_histogram.png
  packet_loss_timeline.png
  report.md
```

## Synthetic trace columns

| Column | Meaning |
|---|---|
| `packet_id` | Monotonic packet identifier. |
| `tx_timestamp_ns` | Synthetic transmit timestamp. |
| `rx_timestamp_ns` | Synthetic receive timestamp. |
| `one_way_delay_ns` | Receiver minus transmitter timestamp. |
| `jitter_ns` | Delay delta relative to previous received packet. |
| `lost` | `0` for received packet, `1` for simulated loss. |
| `burst_id` | Optional impairment burst identifier. |
| `scenario` | Scenario label, for example `baseline`, `jitter_burst`, `loss_burst`. |

## Demo scenarios

| Scenario | Purpose |
|---|---|
| `baseline` | low-jitter reference behavior |
| `jitter_burst` | short interval of increased delay variation |
| `loss_burst` | packet loss cluster |
| `clock_offset` | constant one-way delay bias |
| `mixed_impairments` | combined jitter and loss event |

## SLA summary columns

| Column | Meaning |
|---|---|
| `scenario` | Scenario label. |
| `packets_total` | Total packets generated. |
| `packets_lost` | Number of lost packets. |
| `loss_pct` | Packet loss percentage. |
| `delay_mean_us` | Mean one-way delay in microseconds. |
| `delay_p95_us` | 95th percentile one-way delay. |
| `jitter_p95_us` | 95th percentile absolute jitter. |
| `sla_pass` | `true` / `false` according to configured thresholds. |
| `root_cause_hint` | Simple rule-based explanation. |

## Acceptance criteria

A useful synthetic demo should:

- run without special hardware;
- generate deterministic output with a fixed seed;
- create at least one CSV summary and one report page;
- make jitter and packet loss visible in plots;
- demonstrate how SLA thresholds are evaluated;
- clearly state that synthetic timestamps are not a substitute for hardware timestamping.

## Suggested command

```bash
python tools/generate_synthetic_sla_trace.py --output verification/reports/synthetic_sla_demo/synthetic_trace.csv
python tools/analyze_sla_trace.py \
  --input verification/reports/synthetic_sla_demo/synthetic_trace.csv \
  --output-dir verification/reports/synthetic_sla_demo
```

## Next implementation step

Add the two Python scripts and a generated report template. Keep the first version simple: CSV generation, summary statistics and Markdown output are enough. Plots can be added in the next step.
