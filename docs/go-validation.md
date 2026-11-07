# Go validation evidence

Date: 2026-11-08. Environment: Ubuntu WSL2 x86_64, kernel 6.18.54.1,
Go 1.26.0, GCC, Intel Core i7-8550U @ 1.80 GHz, eight logical CPUs visible.
Windows host: Go 1.27.3 windows/amd64. No production database or hardware fixture
was used. Timings are observations of this development environment.

## Standalone agent stage

Linux `go test -cover ./...`, `go test -race ./...`, and `go vet ./...` passed.
Windows portable tests and vet passed. Linux statement coverage:

| Package | Coverage |
|---|---:|
| cmd/measurement-agent | 78.8% |
| internal/measurement | 88.7% |
| cmd/server-sfp-sla | 1.2% |
| internal/zsocket | 26.3% |

Command: `go test -run='^$' -bench=BenchmarkSynthetic1000 -benchmem -count=3 ./internal/measurement`.
The benchmark includes ID generation, metadata, worker startup, trace generation and
counters, but excludes disk export. The trace contains 1000 modeled packets with
two workers. Three observed samples:

| ns/op | B/op | allocs/op |
|---:|---:|---:|
| 938463 | 76125 | 2878 |
| 795226 | 76093 | 2878 |
| 634293 | 76081 | 2878 |

Substantial variation is expected under WSL2 and host scheduling. These are CPU
work/allocation observations, not network throughput, wire pacing, clock precision
or maximum RSS. No process-level CPU utilization or hardware timing was measured.

A Windows CLI run with defaults and 1000 packets exported a complete package that
the unchanged Python analyzer accepted. Trace SHA-256:
`2092f2ab67e01dce62769ef77c2c7408f1ed06c5530e0cc39f04f116ebbeced4`.
The unit test compares traces across two and seven workers. The analyzer reported
1000 modeled packets, 50 modeled losses, 5% modeled loss, mean synthetic delay
299.970 us, p95 308.875 us, p95 absolute jitter 15.692 us. Default SLA loss threshold
fails intentionally. The analyzer's generic root-cause hint is not a diagnosis of
a physical network and no physical loss burst is implied.

Reproduce with the commands in [the agent guide](go-measurement-agent.md). Generated
full packages are under ignored `bin/agent-results` locally; IDs and wall timestamps
vary. The hash above applies to the deterministic CSV only. Historical evidence
elsewhere in the repository has not been replaced.
