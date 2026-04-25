# Latency model

## One-way delay decomposition

D = D_tx + D_net + D_rx

Where:

- D_tx — transmit path latency
- D_net — network propagation + switching
- D_rx — receive path latency

## With software timestamps

Measured:

D_measured = D + D_os

Where D_os includes:

- scheduler latency
- driver overhead
- interrupt latency

## With FPGA/SFP timestamps

Measured:

D_measured ≈ D_net + constant offset

## Key advantage

The constant offset can be calibrated, leaving pure network delay.

## Calibration strategy

- loopback test
- known cable length
- reference clock comparison
