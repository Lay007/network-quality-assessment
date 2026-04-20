# network-quality-assessment

Hardware-assisted network performance testing system with custom packet protocol, FPGA-based timestamping, and support for RFC 2544 and ITU-T Y.1564 methodologies.

## Overview

This repository contains code, documentation and experiments for practical evaluation of network quality and performance.

The current implementation is based on:
- a Go-based packet generator and receiver
- a custom IPv4 probe packet with an embedded SLA header
- SFP modules with FPGA-based timestamp insertion
- sequence-based loss detection
- delay, one-way delay, jitter and throughput measurements
- RFC 2544 style benchmarking
- ITU-T Y.1564 service-oriented validation

## Documentation

- [Overview](docs/overview.md)
- [Architecture](docs/architecture.md)
- [Packet Format](docs/packet-format.md)
- [Timestamp Model](docs/timestamp-model.md)
- [Metrics](docs/metrics.md)
- [Methodology](docs/methodology.md)
- [Test Modes](docs/test-modes.md)
- [Topology Detection](docs/topology-detection.md)
- [Acceptance Criteria](docs/acceptance-criteria.md)
- [Use Cases](docs/use-cases.md)
- [Results Example](docs/results-example.md)
- [Implementation Notes](docs/implementation-notes.md)
- [Implementation Map](docs/implementation-map.md)
- [Packet Route Diagram](docs/packet-route-diagram.md)

## Examples

- [RFC 2544 Example](examples/rfc2544/readme.md)
- [Y.1564 Example](examples/y1564/readme.md)

## Example Results

See:
- [Sample Test 1 Summary](results/sample-test-1/summary.md)
- [Sample Metrics CSV](results/sample-test-1/metrics.csv)

## Status

- [x] metrics documented
- [x] methodology documented
- [x] packet format documented
- [x] timestamp model documented
- [x] test modes mapped to source files
- [x] topology behavior documented
- [x] example results added
- [x] packet route diagram added
- [ ] add real exported reports from production runs

## Build And Validation

From Linux/WSL:

```bash
make test
make test-cover
make build
make lint
```

Or run the full local pipeline:

```bash
make ci
```

## Runtime Configuration

The service reads database settings from environment variables:

- `SFP_SLA_DB_USER`
- `SFP_SLA_DB_PASSWORD`
- `SFP_SLA_DB_NAME`
- `SFP_SLA_DB_ADDR` (empty for local socket, or `host:port` for TCP)
- `SFP_SLA_VERBOSE=1` to enable verbose logs

See [.env.example](.env.example) and [Setup Guide](docs/setup.md) for full deployment details.
