# Hardware Timestamping Roadmap

This roadmap separates what is already reproducible in software from the hardware proof needed for FPGA/SFP timestamping.

## Current baseline

| Layer | Status | Evidence |
|---|---|---|
| Synthetic trace generation | Ready | `tools/generate_synthetic_sla_trace.py` |
| Offline SLA analysis | Ready | `tools/analyze_sla_trace.py` |
| Report package | Ready | `verification/reports/synthetic_sla_demo/` |
| Trace schema | Documented | `docs/sla-trace-schema.md` |
| FPGA/SFP timestamping | Roadmap | Needs measured hardware package |

## Hardware proof package

A real hardware experiment should contain:

1. topology diagram;
2. clock source and synchronization notes;
3. packet format and timestamp location;
4. capture manifest;
5. raw or reduced trace file;
6. generated plots;
7. summary table;
8. short engineering conclusion;
9. limitations and repeatability notes.

## Measurement stages

| Stage | Goal | Output |
|---:|---|---|
| 1 | Synthetic replay | Stable software report flow |
| 2 | Host timestamp baseline | Show OS scheduling and buffering effects |
| 3 | Hardware timestamp prototype | Show timestamp closer to datapath |
| 4 | Two-node timing run | Estimate one-way delay with clock assumptions stated |
| 5 | SLA report package | Produce a customer-facing report with uncertainty notes |

## Recommended next experiment

Start with a controlled two-port lab setup:

```text
traffic generator -> DUT / link -> capture host
                  -> optional FPGA/SFP timestamp probe
```

Keep packet rate, payload size, link speed, timestamp source and clock assumptions fixed for the first run. The first goal is not maximum performance; it is repeatability and a clean report package.
