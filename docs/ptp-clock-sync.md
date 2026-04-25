# PTP / clock synchronization deep dive

## Why synchronization matters

One-way delay measurement requires both timestamping points to share a common time base.

Without synchronization, the measured delay includes clock offset and drift:

```text
D_measured = D_network + offset_AB + drift_AB(t) + quantization_error
```

Hardware timestamps reduce host noise, but they do not remove clock-domain error by themselves.

## Recommended synchronization options

| Method | Typical role | Notes |
|---|---|---|
| Common reference clock | lab-grade reference | best for controlled benches |
| PPS distribution | simple hardware sync | useful for FPGA-based setups |
| IEEE 1588 PTP | network-distributed time | practical for Ethernet systems |
| Free-running clocks | baseline only | requires drift estimation and compensation |

## Error budget terms

```text
E_sync = E_offset + E_drift + E_jitter + E_asymmetry
```

Where:

- `E_offset` — residual time offset after synchronization
- `E_drift` — frequency mismatch between timestamp clocks
- `E_jitter` — timestamp clock phase noise / short-term instability
- `E_asymmetry` — difference between forward and reverse path delay in sync protocol

## Practical calibration flow

1. Start with a direct loopback path.
2. Measure fixed pipeline latency.
3. Estimate residual timestamp offset.
4. Run a short stability test.
5. Record calibration constants with the report.
6. Repeat calibration after topology, clock, firmware, or SFP changes.

## Acceptance checklist

- timestamp clock frequency is documented
- timestamp resolution is known
- PTP/PPS/reference source is documented
- clock offset is logged before and after run
- drift over the test interval is bounded
- calibration constants are versioned

## Engineering takeaway

Hardware timestamping is only half of the measurement system. The other half is a disciplined time base and a documented calibration procedure.
