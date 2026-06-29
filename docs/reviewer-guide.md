# Reviewer Guide

This guide gives a short path for evaluating `network-quality-assessment`.

## What to review first

| Step | Page | What it proves |
|---:|---|---|
| 1 | [README](../README.md) | Project goal, architecture and runnable synthetic SLA demo |
| 2 | [Synthetic SLA demo](synthetic-sla-demo.md) | End-to-end report generation without FPGA hardware |
| 3 | [SLA trace schema](sla-trace-schema.md) | Input/output contract for reproducible traces |
| 4 | [Measurement credibility](measurement-credibility.md) | Why timestamp location and OS noise matter |
| 5 | [Engineering metrics](engineering-metrics.md) | Delay, jitter, loss and report metrics |
| 6 | [SLA report template](../reports/sla_report.template.md) | Customer-facing report structure |

## Local smoke path

```bash
python tools/generate_synthetic_sla_trace.py --output verification/reports/synthetic_sla_demo/synthetic_trace.csv
python tools/analyze_sla_trace.py --input verification/reports/synthetic_sla_demo/synthetic_trace.csv --output-dir verification/reports/synthetic_sla_demo
```

Expected outputs include a summary CSV, a Markdown report and SVG plots for delay, jitter and packet loss.

## Review checklist

- The synthetic demo can be regenerated from a clean clone.
- The trace schema is stable and documented.
- Report metrics are explained in engineering terms.
- Hardware timestamping claims are separated from synthetic/demo evidence.
- Large generated traces are not required for repository review.

## Current next improvements

1. Add a hardware roadmap from synthetic traces to FPGA/SFP timestamping.
2. Add one compact customer-facing filled SLA report.
3. Keep Go and GitHub Actions dependencies current.
4. Add a comparison table for ping, software timestamps and hardware timestamps.
