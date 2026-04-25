# Traffic analysis: jitter, bursts, SLA violations

## Jitter vs throughput

At increasing load levels:

- queue occupancy increases
- latency variance grows
- jitter becomes a function of buffer dynamics

Non-linear region appears close to line rate.

## Burst detection

Microbursts are short-duration traffic spikes:

- duration: microseconds to milliseconds
- amplitude: near line rate

Detection strategy:

- sliding window packet rate
- inter-arrival time variance
- sudden queue delay increase

## SLA violation patterns

### 1. Random loss
- uniform distribution
- often link quality issue

### 2. Burst loss
- correlated losses
- buffer overflow or congestion

### 3. Delay spikes
- transient congestion
- routing changes

## Correlation analysis

Key idea:

```text
jitter(t) ↔ throughput(t) ↔ queue depth(t)
```

Useful derived metrics:

- jitter vs load curve
- loss vs burst size
- delay percentile (P99, P999)

## Hardware advantage

FPGA timestamps enable:

- microsecond-level burst visibility
- precise inter-packet gap measurement
- correlation of delay with load

## Practical workflow

1. Run controlled traffic profile (ramp / step / burst)
2. Record timestamped packets
3. Build time-series metrics
4. Detect anomalies
5. Map to SLA criteria

## Engineering takeaway

Traffic analysis is only meaningful when the timing reference is precise enough to capture sub-millisecond effects. Hardware timestamps make this possible.
