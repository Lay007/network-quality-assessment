# network-quality-assessment

[![Go](https://img.shields.io/badge/Go-1.22-blue)](#)
[![Build](https://img.shields.io/badge/build-Makefile-informational)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A **hardware-assisted network performance testing system** for practical assessment of delay, jitter, packet loss, and throughput using a custom packet protocol, FPGA-based timestamping, and methodology-oriented test flows.

The repository is centered around real measurement tasks rather than synthetic benchmarks alone. It combines packet generation, custom probe packets, timestamp-aware processing, and structured validation approaches aligned with **RFC 2544** and **ITU-T Y.1564**.

## Highlights

- **Go-based measurement service** and packet-processing pipeline
- **Custom IPv4 probe packet** with embedded SLA-oriented header
- **FPGA/SFP-based timestamp insertion** for hardware-assisted measurements
- **Loss, delay, one-way delay, jitter, and throughput** calculations
- **RFC 2544** style benchmarking support
- **ITU-T Y.1564** service validation support
- deployment scripts and web-console integration for practical use on Linux hosts

## Why this repository matters

Many network-testing repositories focus only on traffic generation or isolated metric collection. `network-quality-assessment` is more interesting because it sits at the boundary between **software measurement logic** and **hardware-assisted timing**.

That makes it useful for:

- network quality assessment projects
- SLA and service-validation tooling discussions
- test-system architecture work involving SFP or FPGA-assisted paths
- practical benchmarking and acceptance workflows for communication links

## Measurement model

At a high level, the system works as follows:

1. the server generates a custom IPv4 measurement packet;
2. the packet traverses the network under test and timestamp-capable SFP/FPGA points;
3. timestamps are inserted on the path and on the return leg;
4. the server receives the probe and derives the relevant metrics.

This model enables a more structured and reproducible view of link behavior than simple host-to-host pings or throughput-only tests.

## Project structure

### Core implementation

- Go measurement service and packet-processing logic
- custom packet format with embedded SLA header
- timestamp-aware metrics engine
- Linux-oriented deployment and runtime scripts
- web-console assets for operational use

### Main documentation

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
- [Setup Guide](docs/setup.md)

### Examples and sample results

- [RFC 2544 Example](examples/rfc2544/readme.md)
- [Y.1564 Example](examples/y1564/readme.md)
- [Sample Test 1 Summary](results/sample-test-1/summary.md)
- [Sample Metrics CSV](results/sample-test-1/metrics.csv)

## Local build and validation

From Linux or WSL:

```bash
make test
make test-cover
make build
make lint
```

Or run the full local validation pipeline:

```bash
make ci
```

The `Makefile` currently includes:

- `build`
- `test`
- `test-cover`
- `vet`
- `gofmt-check`
- `shell-lint`
- `php-lint`
- `web-link-check`
- `lint`
- `ci`

## Binary and runtime notes

The default build target produces:

```text
build/server-sfp-sla
```

The service reads database settings from environment variables, including:

- `SFP_SLA_DB_USER`
- `SFP_SLA_DB_PASSWORD`
- `SFP_SLA_DB_NAME`
- `SFP_SLA_DB_ADDR`
- `SFP_SLA_VERBOSE=1`

See [.env.example](.env.example) and [docs/setup.md](docs/setup.md) for deployment details.

## Deployment model

The Linux deployment flow is script-oriented and intended for practical installation on a host that uses raw sockets and system services.

Typical setup includes:

- dependency installation
- database initialization
- service deployment
- Apache/web-console publishing
- optional user creation for the web interface

## Status

Already documented and represented in the repository:

- metrics model
- methodology mapping
- packet format
- timestamp model
- test modes and topology behavior
- sample results and route diagrams
- setup and deployment guidance

A natural next step would be adding more real exported reports from production or field runs.

## Future improvements

Potential directions for strengthening the repository even further:

- add GitHub Actions for automated validation
- include richer real-world result sets from deployed systems
- expand benchmark and acceptance-report examples
- add deeper operational screenshots or dashboard artifacts
- document more hardware configurations and topology variants

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
