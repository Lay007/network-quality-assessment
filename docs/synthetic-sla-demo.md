# Runnable Synthetic SLA Demo

This page describes the hardware-free demo implemented in this repository. It is intended for reviewers who want to run the measurement and reporting pipeline without FPGA/SFP timestamping hardware.

## What the demo proves

The demo proves the reporting flow, not hardware timestamp accuracy:

```text
synthetic packet trace -> SLA analysis -> summary CSV -> SVG plots -> Markdown report
```

It is useful for CI, documentation, onboarding and portfolio review. Real hardware measurements should use the same report structure but must state timestamp insertion point, clocking, calibration and capture path.

## Run locally

From the repository root:

```bash
python tools/generate_synthetic_sla_trace.py \
  --output verification/reports/synthetic_sla_demo/synthetic_trace.csv

python tools/analyze_sla_trace.py \
  --input verification/reports/synthetic_sla_demo/synthetic_trace.csv \
  --output-dir verification/reports/synthetic_sla_demo
```

Generated outputs:

| Artifact | Purpose |
|---|---|
| `synthetic_trace.csv` | packet-level synthetic transmit/receive timestamp trace |
| `sla_summary.csv` | per-scenario SLA metrics and pass/fail result |
| `report.md` | reviewer-friendly Markdown report |
| `one_way_delay_timeseries.svg` | one-way delay plot |
| `jitter_histogram.svg` | absolute jitter histogram |
| `packet_loss_timeline.svg` | packet-loss event timeline |

## Scenarios

| Scenario | Meaning |
|---|---|
| `baseline` | low-jitter reference behavior |
| `jitter_burst` | short interval of increased delay variation |
| `loss_burst` | clustered packet loss |
| `clock_offset` | constant one-way delay bias |
| `mixed_impairments` | combined delay variation and packet loss |

## SLA thresholds

The analyzer accepts threshold parameters:

```bash
python tools/analyze_sla_trace.py \
  --input verification/reports/synthetic_sla_demo/synthetic_trace.csv \
  --output-dir verification/reports/synthetic_sla_demo \
  --delay-p95-threshold-us 450 \
  --jitter-p95-threshold-us 50 \
  --loss-threshold-pct 1.0
```

A scenario passes when all configured thresholds are met.

## CI workflow

The GitHub Actions workflow `.github/workflows/synthetic-sla-demo.yml` runs the full demo automatically:

1. generate synthetic trace;
2. analyze trace;
3. verify that CSV/SVG/Markdown artifacts are created;
4. upload artifacts from `verification/reports/synthetic_sla_demo/`.

## Limitations

- The demo uses synthetic timestamps.
- It does not estimate FPGA timestamp precision.
- It does not model PHY/MAC implementation details.
- It does not replace calibration or PTP/clocking analysis.
- It exists to make the reporting pipeline reproducible before hardware is attached.

## Next improvements

- Add a compact generated `README.md` inside `verification/reports/synthetic_sla_demo/`.
- Add unit tests for percentile and SLA threshold logic.
- Add CSV schema documentation for real hardware traces.
- Add a real/sanitized trace adapter that writes the same summary format.
