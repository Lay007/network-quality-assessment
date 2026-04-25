# network-quality-assessment

[![Go](https://img.shields.io/badge/Go-1.22-blue)](#)

## 🚀 Hardware timestamping vs software measurements

---

## 🧭 Architecture
![Architecture](docs/assets/architecture.svg)

## 🔁 Packet flow
![Packet Flow](docs/assets/packet_flow.svg)

## 🌐 Topology
![Topology](docs/assets/topology.svg)

## 📊 Benchmark (static)
![Benchmark](docs/assets/benchmark_dashboard.svg)

## 📊 Benchmark (generated from CSV)
![Generated](docs/assets/generated_benchmark.svg)

## 🧪 Case study
👉 docs/case-study.md

---

## 🧠 Hardcore engineering

### Timing error budget
👉 docs/timing-error-budget.md

### Latency model
👉 docs/latency-model.md

---

## 🔬 Generate your own benchmark

```bash
python tools/generate_demo_benchmark.py
```

Outputs:
- results/demo-benchmark/metrics.csv
- docs/assets/generated_benchmark.svg

---

## ⚖️ Measurement approaches

| Approach | Accuracy | Where measured |
|----------|--------|---------------|
| ping | low | OS |
| software | medium | CPU |
| FPGA/SFP | high | datapath |

---

## License
MIT
