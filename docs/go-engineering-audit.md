# Go reliability audit and implementation plan

Reviewed base: `8b99219` (2026-11-08). The working tree was clean. Local and remote
branches were inspected. GitHub listed one open dependency PR (#31) and no open
issues at inspection time. README, CONTRIBUTING, architecture, methodology,
implementation notes, trace schema, SQL schema, Go code and recent history were read.

## Plan and contribution boundaries

1. `refactor/go-core-reliability`, base `master`: correct proven SQL and socket
   defects, bound polling ownership, add regression tests and record limitations.
2. `feat/go-measurement-agent`, base PR 1: independent session engine and synthetic
   CLI with strict configuration and trace-compatible export.
3. `feat/distributed-measurement-bench`, base PR 2: two loopback agents, bounded
   control and probe operations, aggregation and failure tests. Keep hardware
   timestamp input behind an explicit provenance boundary.

No database migration or legacy Ethernet header change is intended. No production
service, database, SNMP device or raw network target was contacted.

## Baseline evidence

Windows: Go 1.27.3 windows/amd64. `go test ./...` and `go vet ./...` failed because
Linux packet socket symbols had no build constraints. `go test -race ./...` failed
because cgo was disabled and no Windows C compiler was available. `gofmt -l .` was
empty. These failures were environment/platform limitations, not failing assertions.

Linux: Ubuntu under WSL2, x86_64, kernel 6.18.54.1, Go 1.26.0, GCC available.
An archive of the unmodified base was tested separately. Tests, race tests and vet
passed; formatting was clean. Statement coverage: legacy command 0.8%, socket
package 11.5%. These tests require neither MySQL nor packet privileges and do not
exercise the physical measurement path. Go was downloaded into a temporary tool
cache outside the repository; no root installation was required.

## Confirmed findings

| Severity | Location | Evidence and disposition |
|---|---|---|
| High | `testBerst.go`, `testLoss.go`, `testDelay.go`, `testY1564.go` | Initial `status=?, datatime=?` received a date then integer 2. Corrected; AST regression checks status argument type/value. |
| High | `testThroughput.go` | Intermediate update supplied nine values for eight placeholders and omitted `rez_64`. Corrected; all literal SQL call sites have an arity regression check. |
| High | legacy SQL error paths | Calling `Close` on rows/DB returned with constructor errors could dereference nil. Removed invalid closes. Pending-query error and cancellation paths are tested with a standard-library fake driver. |
| High | `main.go` polling | Rows were deferred in the process loop before checking errors; six query results were held simultaneously. Polling now materializes and closes each result before launching work, checks iteration errors, limits jobs to one active legacy measurement and suppresses duplicate running keys. |
| High | `internal/zsocket/zsocket.go` | Bind/ring/mmap failures leaked the socket; Close never unmapped the ring. Constructor failure closes its descriptor, concurrent Close is idempotent and waits for active operations. Anonymous mmap and pipe tests need no privileges. |
| High | `internal/zsocket/zsocket.go` | `WriteToBuffer` sliced before checking source length; direct `CopyToBuffer` lacked bounds checks. Reject invalid lengths and callback lengths before publishing a frame; add regression tests. |
| High | `internal/zsocket/zsocket.go` | `poll` received a pointer where the ABI requires a timeout integer. Use typed `unix.Poll` with 50 ms cancellation checks. RX frame bounds and callback ownership are explicit. |
| Medium | `requestSNMP.go` | Defers accumulated a new connection per polling tick; workers could overlap while replacing the module list. Each sample now closes its own connection; context cancellation joins old monitors before replacement. |
| Medium | `testDelay.go` | A new receiver was started every 5 ms on one socket and shared deadlines. Receive serially and stop/close ticker/socket; guard zero-sample division. |
| Medium | `testThroughput.go:genSocket` | Failed generator initialization never notified its owner; zero rate divided by zero. Completion channel always closes, setup failures return counters, rate validation precedes arithmetic. |
| Medium | `config.go` | Driver I/O had no time bounds; DSN construction interpolated database names. Driver config now escapes names and sets 5 s connect/read/write bounds and pool limits. Credentials are never included in new logs. |

The suspected 64-bit `tpacket_req` layout defect was **not confirmed**:
`inet.HOST_INT_SIZE` is four on LP64. Typed kernel helpers replace the handwritten
request calls for portability, not as evidence that the old request layout failed.

## Compatibility and remaining risks

The legacy command and mmap package are explicitly Linux-only. The wire format,
SQL table names, schema and stored numeric field types remain unchanged. Driver
operations now time out after five seconds, pending queries fail closed, and legacy
jobs run serially to avoid competing generators on the same measurement interface.
The startup interface table is refreshed without resetting its auto-increment.

The old hardware functions still contain unchecked SQL operations, blocking retry
loops, float32 timing/statistics, and goroutines outside the new monitor ownership.
Stopping the service cancels polling and joins SNMP monitors; it does **not** promise
a graceful join of every in-flight legacy hardware test. In particular inspect
`TestReal` receiver/watchdog ownership, `testThr` shared counters, ignored raw socket
errors in `testThroughput.go`, and Y.1564 receive workers before production use.
These are explicit residual risks; low command coverage cannot establish race
freedom for unexecuted hardware paths. A hardware fixture and MySQL integration
suite are still needed. The new standalone engine must not call those unsafe paths.

Mmap tests validate user-space lifecycle and bounds, not AF_PACKET ring publication
ordering, NIC behavior, calibrated throughput or FPGA timing. Callbacks must finish
and must copy any borrowed RX bytes retained after return. ARM runtime is pending.
No hardware delay, one-way accuracy or line-rate claim follows from these changes.
