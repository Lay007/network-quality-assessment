# Overview

This repository is dedicated to practical network quality and performance assessment using standardized methodologies and a hardware-assisted measurement path.

## Purpose

The project evaluates communication channel behavior in a structured and repeatable way.

The implementation combines:
- Go-based packet generation
- hardware timestamping via SFP/FPGA modules
- custom measurement packets
- delay, jitter, loss and throughput calculations
- RFC 2544 style benchmarking
- ITU-T Y.1564 service validation

## Practical Model

At a high level:
1. the server generates a custom IPv4 measurement packet
2. the packet is sent toward the first SFP module
3. the module writes timestamps and decides whether to return the packet or forward it further
4. on the return path, additional timestamps are written
5. the server receives the packet and computes metrics
