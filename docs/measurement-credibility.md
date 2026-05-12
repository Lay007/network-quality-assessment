# Measurement Credibility

This project focuses on the difference between host-observed behavior and datapath-observed behavior.

## Why software timestamps are insufficient

Software timestamping is affected by:

- OS scheduling;
- interrupt latency;
- socket buffering;
- driver behavior;
- timestamp placement uncertainty.

As a result, measured delay is often a combination of:

```text
network delay + host processing delay + software jitter
```

## FPGA/SFP timestamping concept

The proposed approach moves timestamps closer to the packet datapath.

```text
packet enters interface
-> hardware timestamp
-> forwarding / processing
-> hardware timestamp
-> exported metrics
```

## Comparison table

| Property | Ping | Software timestamps | FPGA/SFP timestamps |
|---|---|---|---|
| RTT visibility | yes | yes | yes |
| One-way delay | limited | possible | accurate |
| OS noise sensitivity | high | medium | low |
| Timestamp placement | application | socket/driver | datapath |
| SLA confidence | low | medium | high |

## Timing-error budget

Every timing system should document:

- timestamp origin;
- clock synchronization source;
- clock drift assumptions;
- serialization delay;
- queueing visibility;
- reporting granularity.

## Engineering principle

Measurement systems should expose their uncertainty explicitly.

A metric without a documented timing model is not sufficient for engineering-grade SLA validation.
