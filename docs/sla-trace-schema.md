# SLA Trace Schema

This document defines the CSV schema used by the synthetic SLA demo and the recommended schema for future real hardware timestamp captures.

The goal is to keep synthetic, sanitized and hardware-generated traces compatible with the same analysis and reporting pipeline.

## Packet trace schema

| Column | Type | Required | Meaning |
|---|---|---:|---|
| `packet_id` | integer | yes | Monotonic packet identifier. |
| `tx_timestamp_ns` | integer | yes | Transmit timestamp in nanoseconds. |
| `rx_timestamp_ns` | integer / empty | yes | Receive timestamp in nanoseconds. Empty for lost packets. |
| `one_way_delay_ns` | integer / empty | yes | `rx_timestamp_ns - tx_timestamp_ns`. Empty for lost packets. |
| `jitter_ns` | integer / empty | yes | Delay delta relative to the previous received packet. Empty for lost packets. |
| `lost` | integer | yes | `0` for received packet, `1` for lost packet. |
| `burst_id` | integer | recommended | Optional impairment or event identifier. Use `0` for normal packets. |
| `scenario` | string | recommended | Scenario label such as `baseline`, `jitter_burst`, `loss_burst`. |

## Example

```csv
packet_id,tx_timestamp_ns,rx_timestamp_ns,one_way_delay_ns,jitter_ns,lost,burst_id,scenario
0,0,320100,320100,0,0,0,baseline
1,1000000,1320500,320500,400,0,0,baseline
2,2000000,,,,1,2,loss_burst
```

## Summary schema

The analyzer writes `sla_summary.csv` with one row per scenario:

| Column | Meaning |
|---|---|
| `scenario` | Scenario label. |
| `packets_total` | Total packets in the scenario. |
| `packets_lost` | Lost packets in the scenario. |
| `loss_pct` | Packet loss percentage. |
| `delay_mean_us` | Mean one-way delay in microseconds. |
| `delay_p95_us` | 95th percentile one-way delay in microseconds. |
| `jitter_p95_us` | 95th percentile absolute jitter in microseconds. |
| `sla_pass` | `true` if all configured thresholds are met. |
| `root_cause_hint` | Simple rule-based interpretation. |

## Real hardware capture guidance

A hardware timestamp trace should use the same columns whenever possible. Additional columns may be added, but the required columns above should stay stable.

Recommended additional columns for hardware captures:

| Column | Meaning |
|---|---|
| `tx_port` | Transmit port or interface identifier. |
| `rx_port` | Receive port or interface identifier. |
| `clock_domain` | Timestamp clock domain or synchronization source. |
| `ptp_state` | PTP/clock-sync state during capture. |
| `timestamp_point` | Where timestamp was inserted: MAC, PHY, FPGA datapath, driver, OS. |
| `sequence_error` | Optional flag for duplicate/out-of-order packets. |
| `capture_note` | Short note for calibration or known limitations. |

## Compatibility rules

- Use nanoseconds for packet-level timestamps.
- Use microseconds for summary/report values.
- Use empty fields for lost-packet timestamp values.
- Keep `packet_id` monotonic.
- Do not mix host timestamps and FPGA timestamps without a `timestamp_point` note.
- Keep synthetic traces clearly labeled as synthetic.

## Analysis assumptions

The current analyzer assumes:

- one transmit timestamp and one receive timestamp per packet;
- one-way delay is available directly or can be computed before analysis;
- jitter is represented as delay delta between received packets;
- lost packets have empty receive/delay/jitter fields;
- scenario labels can be grouped independently.

These assumptions are intentionally simple so that the same schema can support synthetic CI demos and future sanitized hardware traces.
