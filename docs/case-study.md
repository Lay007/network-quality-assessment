# Case Study: SLA validation on metro Ethernet link

## Scenario

A 1 Gbit/s metro Ethernet link is being validated before service acceptance.

Requirements:
- packet loss <= 0.1 %
- one-way delay <= 1 ms
- jitter <= 100 us
- throughput >= 900 Mbit/s

## Test setup

- probe packets: custom IPv4 SLA format
- rate: 10k packets/sec
- frame size: 256 bytes
- timestamps: FPGA/SFP ingress/egress

Topology:
Host A → SFP/FPGA → Switch → Router → SFP/FPGA → Host B

## Observations

### Without hardware timestamps

- measured delay fluctuates strongly
- jitter inflated by host scheduling
- difficult to distinguish network vs system effects

### With FPGA/SFP timestamps

- delay becomes stable (sub-ms range)
- jitter reflects real network behavior
- packet loss is directly attributable to network path

## Results

| Metric | Target | Measured | Status |
|-------|--------|----------|--------|
| Packet loss | <= 0.10 % | 0.02 % | PASS |
| One-way delay | <= 1 ms | 0.384 ms | PASS |
| Jitter | <= 100 us | 42 us | PASS |
| Throughput | >= 900 Mbit/s | 941 Mbit/s | PASS |

## Engineering insight

The key value is not just measurement, but **measurement with known timestamp origin**.

This allows:

- separating host stack effects from network behavior
- building reproducible SLA validation workflows
- mapping results to RFC 2544 / Y.1564 methodology

## Conclusion

Hardware-assisted timestamping transforms the system from a diagnostic tool into an engineering-grade measurement instrument.