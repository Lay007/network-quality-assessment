# Packet Format

## Packet Structure

Ethernet -> IPv4 -> sfpsla -> optional payload

## IPv4 Header

- custom protocol value: `0x5E`
- source IP: server IP
- destination IP: first SFP module IP
- TTL: `0xFF`

## SLA Header (`sfpsla`)

| Field        | Size    | Description |
|-------------|---------|-------------|
| id          | 1 byte  | packet type / stage identifier |
| dst         | 4 bytes | internal destination / next route |
| merkertime1 | 7 bytes | timestamp T1 |
| merkertime2 | 7 bytes | timestamp T2 |
| merkertime3 | 7 bytes | timestamp T3 |
| number      | 4 bytes | packet sequence number |
| test_type   | 2 bytes | test identifier / mode |
