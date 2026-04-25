# Real-time measurement and anomaly detection pipeline

## Conceptual pipeline

```text
capture → timestamp → decode → metrics → detect → alert → store
```

## Components

### 1. Capture
- raw socket / NIC
- high-rate packet ingestion

### 2. Timestamp
- FPGA/SFP hardware insertion
- nanosecond resolution

### 3. Decode
- custom SLA packet parsing
- sequence tracking

### 4. Metrics
- delay
- jitter
- loss

### 5. Detection
- jitter spikes
- burst loss
- SLA violations

### 6. Alerting
- threshold-based
- event logging
- external integration (syslog / webhook)

## Timing constraints

- processing latency must be << inter-packet gap
- detection window must match SLA requirements

## Engineering challenges

- backpressure under high load
- synchronization between measurement and detection
- timestamp consistency across domains

## Practical insight

Real-time pipeline turns measurement system into **operational monitoring tool**, not just offline analyzer.
