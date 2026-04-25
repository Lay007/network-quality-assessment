# Timing error budget

## Components of measurement error

Total error in delay measurement:

Δt_total = Δt_hw + Δt_sync + Δt_queue + Δt_quant

Where:

- Δt_hw — timestamp insertion latency (FPGA pipeline)
- Δt_sync — clock synchronization error (PTP / local clock drift)
- Δt_queue — buffering and queueing effects
- Δt_quant — timestamp resolution (clock granularity)

## Hardware timestamp path

Typical contributions:

- FPGA insertion latency: 10–100 ns
- PHY latency variation: < 50 ns
- SFP processing: deterministic

## Software timestamp path

Typical contributions:

- kernel scheduling: 10–1000 us
- NIC queueing: 10–500 us
- interrupt coalescing: 50–300 us

## Key insight

Hardware timestamps reduce dominant error terms by **3–4 orders of magnitude**.

## Engineering implication

This enables:

- accurate one-way delay measurement
- meaningful jitter analysis
- SLA validation instead of approximate diagnostics
