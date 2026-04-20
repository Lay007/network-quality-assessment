# Packet Route Diagram

## Measurement Route

```mermaid
sequenceDiagram
    participant G as Go Generator / Server
    participant S1 as SFP Module 1 (FPGA)
    participant S2 as SFP Module 2 / Remote Path
    participant R as Server Receiver

    G->>S1: IPv4 + sfpsla probe packet
    Note right of S1: Insert T1
    alt forward path
        S1->>S2: forward packet
        Note right of S2: Insert T2
        S2->>S1: return packet
        Note right of S1: Insert T3
        S1->>R: final packet to server
    else reflect locally
        S1->>R: reflected packet with timestamps
    end
```

## Interpretation

- T1: first module observation
- T2: remote/intermediate observation
- T3: return observation
- sequence number supports packet loss estimation
- sampled packets can be timestamped while broader traffic still supports loss statistics
