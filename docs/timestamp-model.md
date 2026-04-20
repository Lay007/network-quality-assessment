# Timestamp Model

## Conversion to Microseconds

The current implementation converts delay values using:

`delay_us = ticks * 1e6 / 2^32`

## Delay Calculations

- Two-way delay may be derived from `T3 - T1`
- Another mode may derive delay from `now - T2`
- One-way components are estimated from `T2 - T1` and `T3 - T2`

## Jitter

Jitter is computed as variation between consecutive delay measurements.
