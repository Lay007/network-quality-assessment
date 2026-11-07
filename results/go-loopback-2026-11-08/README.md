# Go loopback software evidence: 2026-11-08

Captured from clean source revision `7bed4448987c6b331e414b78c4b2b90706b0d659`
with Go 1.26.0 on Ubuntu WSL2 x86_64. This package contains three real local
software bench sessions, 100 probes per agent per run (600 completed observations,
zero missing authenticated replies). It is not FPGA, one-way, physical two-machine
or capacity evidence. Configured pacing deliberately limits traffic.

- `validation.json`: distributions, child CPU and maximum individual child RSS.
- `checks.json`: test counts, coverage, builds and allocation benchmark samples.
- `packages/`: complete aggregate JSON and per-agent metadata/RTT CSVs.
- `SHA256SUMS`: artifact hashes; verify from this directory with `sha256sum -c SHA256SUMS`.

Reproduce from the repository root:

```bash
go build -o bin/measurement-agent ./cmd/measurement-agent
go build -o bin/measurement-bench ./cmd/measurement-bench
python3 scripts/run_go_bench.py --runs 3 --packets 100
```

Each run gets fresh IDs, nonce values, wall timestamps and timing observations.
The saved evidence is an observation rather than a golden performance threshold.
The runner's RTT percentiles use Python `statistics.quantiles` with inclusive
interpolation. CPU totals cover the two agents and three controller processes;
maximum RSS is the largest individual child peak, not the sum. Full methodology
and limits are in [the bench guide](../../docs/distributed-go-bench.md).
