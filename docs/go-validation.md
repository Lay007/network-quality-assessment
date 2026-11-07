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

## Final distributed bench stage

Source revision: `7bed4448987c6b331e414b78c4b2b90706b0d659` (clean build).
[The committed result package](../results/go-loopback-2026-11-08/README.md) contains
all three aggregates, per-agent RTT CSV/JSON, hashes and machine-readable checks.

All **45 tests** passed on Linux. Linux race tests and vet passed; formatting and
both repository link checkers passed. Windows portable tests and vet passed.
The concurrency/failure suites also passed ten repeated runs before the final build.

| Package | Tests | Final statement coverage |
|---|---:|---:|
| cmd/measurement-agent | 4 | 82.6% |
| cmd/measurement-bench | 2 | 79.1% |
| internal/measurement | 8 | 76.8% |
| internal/bench | 9 | 79.7% |
| cmd/server-sfp-sla | 13 | 1.2% |
| internal/zsocket | 9 | 26.3% |

Whole-repository Linux AMD64 and ARM64 builds passed. Standalone agent/controller
ARMv7 builds passed (`GOARM=7`). Whole-repository ARMv7 remains blocked by missing
host-size definitions in the legacy socket dependency. ARM runtime, a physical
two-machine experiment, MySQL/hardware integration and calibrated FPGA timestamps
were not tested.

The process-level runner launched two built agents and three successive controller
runs. Both agents exited cleanly after SIGTERM. Each run observed 200 authenticated
loopback replies out of 200 probes (600/600 in total).

| Run | Mean RTT (us) | p95 RTT (us) | p99 RTT (us) | Max RTT (us) |
|---:|---:|---:|---:|---:|
| 1 | 363.519 | 491.450 | 729.439 | 2396.373 |
| 2 | 261.785 | 407.456 | 475.057 | 633.973 |
| 3 | 332.280 | 505.093 | 810.960 | 2029.631 |

These are application send-to-reply software observations on one WSL2 host.
They include OS scheduling and echo work and are not precision specifications.
The complete runner took 5.588 s. Child CPU totals were 0.360 s user and 0.773 s
system. Maximum individual child RSS was 19048 KiB, not aggregate memory.
Combined reply payload rates were 89.210, 91.476 and 91.438 kbit/s over controller
session windows (received replies * 64 bytes * 8 / elapsed seconds). This deliberately
paced payload rate excludes request traffic, protocol overhead and startup; it does
not measure network capacity. The raw aggregate start/end times permit recomputation.

The final synthetic benchmark (same command as above) measured 1139923, 893568 and
1322150 ns/op; 92693, 92667 and 92644 B/op; 2878 allocs/op in each sample. Sample
representation now also accommodates software RTT fields, so these are a separate
stage observation rather than a controlled comparison with the earlier benchmark.
A clean Linux synthetic CLI run reproduced the exact same CSV SHA-256 as the
Windows run, and the unchanged Python analyzer generated its report successfully.

## Review status and next work

All three changes are prepared as dependent drafts rather than merged releases.
The standalone model and local bench are independently reproducible; the core draft
is a bounded reliability improvement, **not completion of every legacy lifecycle
and shared-state repair**. The residual hardware worker risks listed in
[the audit](go-engineering-audit.md) remain before production qualification.
Prioritize context ownership for those workers and a MySQL/socket fixture, then
confirm the FPGA timestamp contract and run an authorized remote/hardware experiment.


A subsequent review added concurrent cancellation requests within the existing
cleanup deadline, with a regression where one cancel endpoint stalls until the
other agent is canceled. The current suite has **46 tests**; Linux tests/race/vet
and Windows portable tests/vet pass after that change. The archived 45-test and
performance snapshot above remains tied to its explicitly recorded source revision.
