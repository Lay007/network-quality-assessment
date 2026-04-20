# Implementation Map

- `cmd/server-sfp-sla/main.go` — orchestration, packet construction, real SLA receive logic
- `cmd/server-sfp-sla/struct.go` — packet header and test structs
- `cmd/server-sfp-sla/testDelay.go` — delay and jitter
- `cmd/server-sfp-sla/testLoss.go` — loss
- `cmd/server-sfp-sla/testThroughput.go` — throughput
- `cmd/server-sfp-sla/testBerst.go` — burst/back-to-back style testing
- `cmd/server-sfp-sla/testY1564.go` — Y.1564 tests
- `cmd/server-sfp-sla/requestSNMP.go` — SNMP monitoring
