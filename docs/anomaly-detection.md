# Anomaly detection for SLA validation

## Types of anomalies

### 1. Jitter spikes

Condition:

```text
jitter(t) > mean(jitter) + k * std(jitter)
```

Typical k = 3

---

### 2. Burst loss

Condition:

```text
loss(t..t+n) > threshold AND contiguous
```

Detects microburst congestion.

---

### 3. Delay excursions

Condition:

```text
delay(t) > SLA_limit
```

---

## Sliding window detector

Window size:
- 10 ms to 1 s depending on link

Metrics:
- mean
- variance
- max

---

## Example output

```text
[ALERT] jitter spike at t=23.4s: 120 us
[ALERT] burst loss detected: 12 packets lost
[ALERT] delay violation: 1.8 ms > 1 ms
```

---

## Engineering insight

Without hardware timestamps:
- anomalies blurred

With FPGA timestamps:
- anomalies are sharply defined

---

## SLA scoring

Composite score:

```text
score = w1*loss + w2*jitter + w3*delay
```

Used for pass/fail automation.
