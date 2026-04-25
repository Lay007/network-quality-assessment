# Demo: hardware-assisted SLA measurement flow

This demo is a reproducible documentation scenario that shows the expected workflow of the project.

## Goal

Show the complete chain:

```text
probe generation -> timestamped packet path -> metrics -> anomaly detection -> report
```

## Run

From the repository root:

```bash
python tools/generate_demo_benchmark.py
```

Generated artifacts:

```text
results/demo-benchmark/metrics.csv
docs/assets/generated_benchmark.svg
```

## Scenario

- service type: 1G Ethernet SLA validation
- packet type: custom IPv4 SLA probe
- timestamp source: FPGA/SFP datapath
- test duration: 60 seconds
- acceptance model: delay / jitter / loss / throughput

## Expected interpretation

The generated dashboard demonstrates:

- stable sub-ms one-way delay
- bounded jitter
- rare loss events
- throughput close to line rate

## Engineering value

The demo is intentionally simple, but it keeps the structure of a real measurement workflow:

1. deterministic dataset generation;
2. report-oriented metric extraction;
3. visualization suitable for README and CI-generated artifacts;
4. clear separation between synthetic demo data and production measurements.
