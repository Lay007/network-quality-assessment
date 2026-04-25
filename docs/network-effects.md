# Network effects under load

## Key phenomena

### 1. Jitter vs load

As utilization approaches line rate:

- queue depth increases
- jitter grows non-linearly

### 2. Microbursts

Short bursts cause:

- temporary queue overflow
- packet loss spikes

### 3. Bufferbloat

Large buffers:

- reduce loss
- increase latency dramatically

## Measurement advantage

With FPGA timestamps:

- jitter reflects actual queue dynamics
- delay spikes are visible at microsecond scale

## Without hardware timestamps

- jitter inflated by OS noise
- microbursts masked

## Engineering takeaway

Accurate SLA validation requires:

- timestamp close to datapath
- controlled traffic profile
- correlation between load and delay
