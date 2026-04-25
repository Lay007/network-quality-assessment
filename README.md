# network-quality-assessment

[![Go](https://img.shields.io/badge/Go-1.22-blue)](#)
[![Build](https://img.shields.io/badge/build-Makefile-informational)](#)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 🚀 Hardware timestamping vs software measurements

A **hardware-assisted network performance testing system** combining:
- custom IPv4 probe packets
- FPGA/SFP timestamping
- SLA-oriented metrics (RFC 2544 / ITU-T Y.1564)

---

## 🧭 Architecture

![Architecture](docs/assets/architecture.svg)

---

## 🔁 Packet flow

![Packet Flow](docs/assets/packet_flow.svg)

---

## 🌐 Test topology

![Topology](docs/assets/topology.svg)

---

## 📊 Example report

![Report](docs/assets/report_example.svg)

---

## ⚖️ Measurement approaches

| Approach | Accuracy | Where measured | Use case |
|----------|--------|---------------|----------|
| ping | low | OS | connectivity |
| software timestamps | medium | CPU | lab tests |
| FPGA/SFP timestamps | high | datapath | SLA / engineering |

---

## 💡 Why this is interesting

Unlike typical tools, this system:

- separates **network delay vs host delay**
- embeds measurement data directly into packets
- enables **engineering-grade SLA validation**

---

## 🔧 Build

```bash
make build
make test
```

---

## 📁 Docs

See `/docs` for full methodology and implementation details.

---

## 📜 License

MIT
