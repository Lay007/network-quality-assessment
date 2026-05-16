# Experiment Manifests

This directory contains measurement scenarios for reproducible network-quality analysis.

## Available manifests

| Manifest | Purpose |
|---|---|
| `sla_demo.yaml` | baseline hardware timestamping SLA demonstration |
| `software_vs_hardware_timestamp.yaml` | comparison between host-observed and datapath-observed timing |

## Manifest purpose

A measurement manifest should document:

- topology;
- timestamp origin;
- traffic profile;
- expected metrics;
- report outputs;
- timing assumptions.

## Recommended workflow

```text
select scenario
-> generate or capture traffic
-> extract timing metrics
-> build SLA report
-> review timestamp confidence
```

## Future scenarios

- PTP clock drift analysis;
- queue buildup detection;
- packet-loss burst classification;
- multi-hop SLA validation.
