# Real vs Synthetic Validation Guide

This guide explains how the synthetic SLA demo should be connected to future real measurement runs.

The synthetic demo is useful because it is reproducible and hardware-free. Real measurements are useful because they expose clock, buffering, queueing, and deployment effects that a synthetic trace can only approximate.

## Validation goal

A mature measurement project should show both:

```text
synthetic trace -> known expected behavior -> analysis pipeline check
real trace      -> measured deployment behavior -> engineering conclusion
```

Synthetic data validates the toolchain. Real data validates the measurement method.

## What synthetic traces prove

Synthetic traces are good for:

- checking parser correctness;
- checking report generation;
- testing threshold rules;
- creating reproducible CI examples;
- demonstrating packet loss and jitter visualization;
- keeping the repository reviewable without hardware.

Synthetic traces do not prove actual deployment accuracy.

## What real traces prove

Real traces are needed for:

- clock synchronization assessment;
- host-vs-datapath timestamp comparison;
- queueing and burst behavior;
- packet capture limitations;
- link-specific jitter and delay patterns;
- practical SLA conclusions.

## Comparison workflow

| Step | Synthetic trace | Real trace |
|---|---|---|
| Input contract | generated CSV | captured CSV with metadata |
| Expected behavior | known by construction | estimated from setup notes |
| Analysis command | same toolchain | same toolchain |
| Output | summary, plots, report | summary, plots, report |
| Review question | does the toolchain work? | is the measurement credible? |

## Minimum metadata for real traces

A real trace should be accompanied by a metadata file:

```yaml
trace_id: "YYYYMMDD_site_test"
measurement_type: "software"  # software / hardware / mixed
topology: ""
clock_source: ""
sync_method: ""
packet_rate_pps: null
packet_size_bytes: null
capture_duration_s: null
link_description: ""
operator: ""
notes: ""
```

## Acceptance criteria

A real-vs-synthetic comparison is useful when:

- both traces follow the same schema;
- both are processed by the same analysis command;
- summary metrics are reported in the same units;
- plots use comparable axes where practical;
- the report clearly separates toolchain validation from deployment conclusions.

## Suggested report sections

1. Measurement purpose.
2. Synthetic trace configuration.
3. Real trace setup and metadata.
4. Common analysis command.
5. Metric comparison table.
6. Plot comparison.
7. Interpretation.
8. Limitations.
9. Next measurement action.

## Practical decision rule

Use synthetic traces to keep CI and documentation honest. Use real traces to make engineering claims about networks.
