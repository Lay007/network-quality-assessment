# network-quality-assessment

[![Docs Check](https://github.com/Lay007/network-quality-assessment/actions/workflows/docs-link-check.yml/badge.svg)](https://github.com/Lay007/network-quality-assessment/actions/workflows/docs-link-check.yml)

A **hardware-assisted network measurement and SLA validation system** combining FPGA timestamping, packet-level telemetry and real-time analytics.

---

## 🚀 Quick navigation

- [Architecture](#-architecture)
- [Packet flow](#-packet-flow)
- [Benchmark & analytics](#-benchmark--analytics)
- [Case study](#-case-study)
- [Network analysis](#-network-analysis)
- [Real-time system](#-real-time-system)
- [Root cause analysis](#-root-cause-analysis)
- [Clock synchronization](#-clock-synchronization)
- [Hardcore engineering](#-hardcore-engineering)
- [Quick start](#-quick-start)
- [Demo](#-demo)

---

## 🧭 Architecture

![Architecture](docs/assets/architecture.svg)

---

## 🔁 Packet flow

![Packet Flow](docs/assets/packet_flow.svg)

---

## 🌐 Topology

![Topology](docs/assets/topology.svg)

---

## ⚙️ FPGA timestamp pipeline

![FPGA](docs/assets/fpga_timestamp_pipeline.svg)

---

## 📦 Packet decode (Wireshark-style)

![Packet](docs/assets/packet_format_decode.svg)

---

## 📊 Benchmark & analytics

![Benchmark](docs/assets/benchmark_dashboard.svg)
![Generated](docs/assets/generated_benchmark.svg)
![Analytics](docs/assets/analytics_correlation.svg)

---

## 🧪 Case study

👉 [Metro Ethernet SLA validation](docs/case-study.md)

---

## 🌐 Network analysis

- 👉 [Network effects](docs/network-effects.md)
- 👉 [Traffic analysis](docs/traffic-analysis.md)
- 👉 [Anomaly detection](docs/anomaly-detection.md)

---

## 🧠 Root cause analysis

👉 [Root cause analysis rules](docs/root-cause-analysis.md)

---

## ⚡ Real-time system

👉 [Real-time pipeline](docs/realtime-pipeline.md)

---

## 📡 Monitoring integration

👉 [Monitoring integration](docs/integration-monitoring.md)

---

## ⏱️ Clock synchronization

👉 [PTP / clock sync](docs/ptp-clock-sync.md)

---

## 🧠 Hardcore engineering

- 👉 [Timing error budget](docs/timing-error-budget.md)
- 👉 [Latency model](docs/latency-model.md)

---

## 🔬 Quick start

Generate demo dataset and graphs:

```bash
python tools/generate_demo_benchmark.py
```

---

## 🎯 Demo

👉 [Run demo scenario](examples/demo/README.md)

This demo shows a full measurement flow:

```text
probe -> timestamp -> metrics -> detection -> report
```

---

## 📜 License

MIT
