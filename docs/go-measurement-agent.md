# Standalone Go measurement agent

Build with Go 1.26 or newer, from the repository root:

```bash
go build -o bin/measurement-agent ./cmd/measurement-agent
bin/measurement-agent --packets 1000 --seed 1 --output bin/agent-results
```

On Windows use `bin/measurement-agent.exe`. The command prints a new session
directory containing `metadata.json` and `trace.csv`. No database, root, raw socket,
FPGA or network target is required. Each run gets a cryptographically random ID.
The packet trace is deterministic for a fixed model configuration even if worker
count changes; IDs, wall-clock execution times and software version are not.

The metadata includes configuration version, full configuration and SHA-256,
software revision (or development), Go version, platform, start/end UTC times,
status, structured errors, modeled packet counts/loss and timestamp provenance.
The only supported mode in this stage is synthetic. The interface is a model label,
not a physical NIC. Modeled timestamps start at zero; they are not host timestamps.

## Analysis and result semantics

Reuse the existing Python analyzer on the printed directory:

```bash
python tools/analyze_sla_trace.py --input bin/agent-results/SESSION_ID/trace.csv \
  --output-dir bin/agent-results/SESSION_ID/analysis
```

Replace `SESSION_ID` with the directory printed by the agent. The CSV follows
[sla-trace-schema.md](sla-trace-schema.md) without extensions. Lost packets have
empty RX/delay/jitter fields; jitter spans successive received packets. The default
model loses every twentieth packet (5% modeled loss), so the analyzer's default
1% loss SLA fails intentionally. Use `--loss-every 0` for a no-loss model. This is
not observed network loss. Python remains responsible for percentile analysis and
plots; the Go engine exports the trace and basic session accounting.

## Configuration and lifecycle

Use `--help` for flags. `--config examples/measurement-agent-config.json` loads a
JSON configuration; duration fields are integer nanoseconds. Unknown fields and
multiple objects are rejected. Files are limited to 64 KiB. A file can be combined
with `--output`, but not measurement flags. Omitted JSON fields inherit defaults.

Limits: 1..10000 packets, 1..16 workers, 1 us..1 s modeled interval, 1 ms..1 min
computation timeout, nonnegative modeled delay up to one second, and variation no
greater than the delay. These keep memory and work bounded. Synthetic computation
does not sleep for modeled packet intervals. A ten-second modeled trace can be
computed in milliseconds; its speed does not describe packet throughput.

SIGINT/SIGTERM, parent context cancellation and `--timeout` stop computation and
join all workers. A canceled/timed-out session exports its completed contiguous
prefix and error metadata, then exits nonzero. Uncomputed packets are not counted
as loss. Invalid configuration fails before starting a session. Export uses a
private temporary directory and publishes after all files close successfully;
existing session packages are rejected. Interrupted disk writes can leave an
unpublished `.session-*` directory; those are not completed packages.

## Validation and limits

```bash
go test ./...
go vet ./...
# Linux with GCC/cgo:
go test -race ./...
go test -run='^$' -bench=BenchmarkSynthetic1000 -benchmem -count=3 ./internal/measurement
```

The legacy command and AF_PACKET package are Linux-only. Portable session and CLI
packages also run on Windows. Tests verify deterministic output across worker
counts, bounds, loss/jitter semantics, cancellation, strict configuration and
non-overwriting exports. See [the reliability audit](go-engineering-audit.md) for
legacy risks. Synthetic evidence establishes reproducibility of the model and
report contract, not timing accuracy, network capacity or hardware performance.
