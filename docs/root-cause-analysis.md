# Root cause analysis rules

## Goal

Move from raw metric violations to engineering explanations.

```text
metric anomaly → correlation → probable root cause → recommended action
```

## Rule examples

| Symptoms | Probable cause | Recommended action |
|---|---|---|
| jitter increases with throughput, loss remains low | queue buildup / congestion | inspect QoS, shaping and buffer configuration |
| burst loss with short delay spike | microburst overflow | check bursty source, policing, switch buffers |
| delay increases but jitter remains low | path change or added processing stage | verify route, tunnel, firewall or shaping node |
| random isolated loss | physical layer / optical issue | inspect SFP, fiber, counters, CRC/FEC |
| one-way delay asymmetry | path asymmetry or sync issue | compare forward/reverse paths and PTP status |

## Decision tree

```text
loss?
├─ bursty → buffer overflow / microburst
├─ random → physical layer / optical / CRC issue
└─ none
   └─ jitter?
      ├─ grows with load → congestion / queueing
      └─ isolated spikes → transient scheduling or route event
```

## Confidence score

Each hypothesis should include:

- matched symptoms
- supporting metrics
- missing evidence
- confidence level

## Example

```text
Root cause: queue buildup under load
Evidence:
- jitter grows after 900 Mbit/s
- delay P99 increases
- packet loss remains low
Confidence: 0.82
Recommendation: check QoS policy and egress queue configuration
```
