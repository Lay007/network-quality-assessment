# Sample Test Report

## Scenario
Baseline network test with two SFP modules in a serial path.

## Topology
Server -> SFP1 -> SFP2 -> return path -> Server

## Metrics

| Metric             | Value   |
|-------------------|---------|
| RTT avg           | 12.4 us |
| RTT max           | 16.8 us |
| RTT min           | 11.7 us |
| OWD forward       | 6.1 us  |
| OWD reverse       | 6.3 us  |
| Jitter            | 0.4 us  |
| Packet loss       | 0.0001  |
| Throughput        | 9.8 Gbps |

## Observations
- stable two-way delay
- low jitter
- negligible packet loss
- no abnormal qualitative behavior in the example scenario

## Conclusion
The path is operating within expected parameters for the sample test case.
