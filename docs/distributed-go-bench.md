# Distributed Go measurement bench

The bench coordinates **two local agents on Linux loopback**. It validates session
control, cancellation, result aggregation and application-level software RTT.
It does not establish two-machine behavior, synchronized starts, one-way delay,
FPGA timing or network capacity.

## Build and run the complete scenario

Go 1.26+, Python 3 and an unprivileged Linux account are sufficient:

```bash
go build -o bin/measurement-agent ./cmd/measurement-agent
go build -o bin/measurement-bench ./cmd/measurement-bench
python3 scripts/run_go_bench.py --runs 3 --packets 100
```

The runner creates a random token in child environments, launches two agent
processes on ephemeral loopback TCP/UDP ports, runs three controller sessions,
exports full packages and terminates both agents. Inspect
`bin/local-bench/validation.json` for actual RTT distributions, loss, execution
time, child CPU time and maximum individual child RSS where available. No token
is written to results. The runner sends at most 200 probes per agent per run and
at most ten runs. On Windows, portable tests work; forced process termination in
the Python runner is not evidence of graceful signal shutdown.

## Manual Linux deployment

Use three terminals in the same network namespace. Set the **same** token in all
three environments, generated locally, for example with
`export NQA_BENCH_TOKEN="$(python3 -c 'import secrets; print(secrets.token_hex(24))')"`
in the first terminal, then transfer its value privately to the other terminals.
Do not store the token in Git or command-line arguments.

```bash
# Terminal A
bin/measurement-agent serve --id local-a --listen 127.0.0.1:8851 --probe 127.0.0.1:9851
# Terminal B
bin/measurement-agent serve --id local-b --listen 127.0.0.1:8852 --probe 127.0.0.1:9852
# Terminal C
bin/measurement-bench --agents http://127.0.0.1:8851,http://127.0.0.1:8852 \
  --packets 20 --output bin/bench-results
```

`--mode synthetic` on the controller coordinates deterministic model sessions.
The printed package contains `aggregate.json` and `agents/MEASUREMENT_ID/` with
per-agent metadata and CSV. Software runs use `rtt.csv`; model runs use the existing
`trace.csv` contract. Aggregate JSON maps agent identity, control session IDs and
measurement session IDs. It retains the last observed state on control failure;
that state does not imply the agent continued after the controller's cleanup.

For a standalone software session against a running echo agent:

```bash
bin/measurement-agent --mode udp --target 127.0.0.1:9852 --workers 1 \
  --interval 10ms --packets 20 --output bin/agent-results
```

## Lifecycle, bounds and security

HTTP requests require `Authorization: Bearer TOKEN` (`NQA_BENCH_TOKEN`, 16..4096
bytes). GET `/v1/info` returns identity, platform, software, UDP endpoint and clock
provenance. POST `/v1/sessions` takes `control_session_id` (32 lowercase hex digits)
and the agent configuration; it returns 202 and a running record. GET
`/v1/sessions/ID` returns state/results; DELETE requests cancellation. Duplicate IDs
return 409. A controller generates and remembers IDs **before** POST, so it can
request cleanup even if the acceptance response is lost. There are no automatic
start retries or claims of exactly-once execution across restarts.

States: running -> stopping (optional) -> completed/canceled/timeout/network/failed.
Sessions belong to the agent lifetime rather than the HTTP request. Signal shutdown
cancels probes and joins workers. Each agent allows two active sessions and retains
at most sixteen records, evicting the oldest terminal record. Sessions are in-memory;
results are not durable until fetched and exported. An agent crash can lose them.

All bind/target/control addresses must be numeric loopback endpoints. External
addresses, DNS names, URL credentials, redirects and proxies are rejected/disabled.
HTTP time bounds: 2 s client requests, 1 s dial, 2 s header read, 3 s server read/write,
5 s idle. Start bodies are limited to 64 KiB and responses to 8 MiB. Controller
execution is bounded by 30 s; best-effort cancellation uses a separate 2 s cleanup
context. A dead agent may make cleanup impossible; the agent's own session timeout
still bounds work if it remains alive. No public listener, TLS deployment or real
remote-machine deployment is claimed by this local-only milestone.

Software sessions use one connection and one outstanding UDP probe, <=1000 packets,
>=10 ms pause after a completed probe, and a 1..200 ms per-probe deadline. A 64-byte
probe contains magic/version, random measurement ID, sequence, random nonce and a
truncated 128-bit HMAC-SHA256. Echo verifies the MAC, size and loopback source before
replying with identical bytes. Responses must match the complete outstanding probe;
stale or unrelated datagrams are ignored until its original deadline. This is a
new bench protocol; the legacy SFP Ethernet protocol is unchanged.

## Timing and loss interpretation

RTT is measured with `time.Time.Sub` retaining Go's local monotonic component,
from immediately before the UDP write to receipt of its matching echo. It includes
application scheduling, queues, kernel work and echo processing. TX/RX Unix nanoseconds
are diagnostics and may change with wall-clock adjustments; RTT does not subtract
those Unix values. No remote timestamp is used, so clock synchronization is not
required for RTT. RTT delta is a difference between successive received RTTs;
it is not one-way packet-delay variation.

`packets_sent` counts successful UDP writes. `packets_attempted` counts completed
observations; timeout without an echo counts as loss, while a canceled in-flight
probe is `packets_unfinished` and not loss. Network errors end the session with an
error, rather than fabricating a complete loss trace. Loss means no authenticated
reply before the configured deadline; it does not identify whether request, reply,
scheduling or endpoint failure caused the missing observation.

Synthetic traces remain compatible with the existing Python analyzer. **Do not feed
`rtt.csv` to that analyzer**: it assumes one-way delay. See
[the schema distinction](sla-trace-schema.md#software-round-trip-results).

`measurement.TimestampInput` and `TimestampEvent` define the future capture boundary:
integer timestamp, packet ID, clock domain, capture point, provenance, cancellation
and close. No hardware provider or FPGA register interface is implemented. A future
provider needs a confirmed hardware contract and independent synchronization,
calibration and timestamp-point evidence before one-way results can be exported.

## Platforms and verification

```bash
go test ./...
go test -race ./...
go vet ./...
GOOS=linux GOARCH=amd64 go build ./...
GOOS=linux GOARCH=arm64 go build ./...
GOOS=linux GOARCH=arm GOARM=7 go build ./cmd/measurement-agent ./cmd/measurement-bench
```

AMD64 Linux runtime is validated under WSL2. ARM64 whole-repository and ARMv7
standalone command builds pass cross-compilation; ARM hardware execution is pending.
Whole-repository ARMv7 is blocked by the legacy `newtools/zsocket/inet` dependency,
which lacks Linux 32-bit ARM host-size definitions. New commands do not import it.
A real two-machine test remains pending; extending to remote targets requires an
explicit transport/security configuration and authorization rather than removing
the loopback gate silently.
