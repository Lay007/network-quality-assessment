#!/usr/bin/env python3
"""Generate a deterministic synthetic SLA packet trace.

The script intentionally has no third-party dependencies so it can run in CI and
on a clean workstation. It does not model hardware timestamping accuracy; it
creates a public, reproducible trace that exercises the reporting pipeline.
"""

from __future__ import annotations

import argparse
import csv
import random
from pathlib import Path


NS_PER_MS = 1_000_000
NS_PER_US = 1_000


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Generate a synthetic SLA packet trace")
    parser.add_argument("--output", required=True, help="Output CSV path")
    parser.add_argument("--packets", type=int, default=2500, help="Number of packets to generate")
    parser.add_argument("--seed", type=int, default=42, help="Random seed")
    parser.add_argument("--period-us", type=float, default=1000.0, help="Packet period in microseconds")
    parser.add_argument("--base-delay-us", type=float, default=320.0, help="Baseline one-way delay")
    parser.add_argument("--baseline-jitter-us", type=float, default=3.0, help="Baseline jitter sigma")
    return parser.parse_args()


def scenario_for_packet(packet_id: int, total_packets: int) -> tuple[str, int]:
    """Return scenario label and burst id for a packet index."""
    q1 = total_packets // 5
    q2 = 2 * total_packets // 5
    q3 = 3 * total_packets // 5
    q4 = 4 * total_packets // 5

    if q1 <= packet_id < q1 + max(40, total_packets // 25):
        return "jitter_burst", 1
    if q2 <= packet_id < q2 + max(30, total_packets // 30):
        return "loss_burst", 2
    if q3 <= packet_id < q3 + max(80, total_packets // 12):
        return "clock_offset", 3
    if q4 <= packet_id < q4 + max(90, total_packets // 10):
        return "mixed_impairments", 4
    return "baseline", 0


def delay_for_scenario(
    rng: random.Random,
    scenario: str,
    base_delay_ns: int,
    baseline_jitter_ns: int,
) -> int:
    if scenario == "jitter_burst":
        return base_delay_ns + int(rng.gauss(0.0, 35 * NS_PER_US))
    if scenario == "clock_offset":
        return base_delay_ns + 85 * NS_PER_US + int(rng.gauss(0.0, baseline_jitter_ns))
    if scenario == "mixed_impairments":
        return base_delay_ns + 45 * NS_PER_US + int(rng.gauss(0.0, 25 * NS_PER_US))
    return base_delay_ns + int(rng.gauss(0.0, baseline_jitter_ns))


def is_lost(rng: random.Random, scenario: str) -> bool:
    if scenario == "loss_burst":
        return rng.random() < 0.18
    if scenario == "mixed_impairments":
        return rng.random() < 0.06
    return rng.random() < 0.001


def generate_rows(args: argparse.Namespace) -> list[dict[str, object]]:
    rng = random.Random(args.seed)
    period_ns = int(args.period_us * NS_PER_US)
    base_delay_ns = int(args.base_delay_us * NS_PER_US)
    baseline_jitter_ns = int(args.baseline_jitter_us * NS_PER_US)

    rows: list[dict[str, object]] = []
    previous_received_delay: int | None = None

    for packet_id in range(args.packets):
        scenario, burst_id = scenario_for_packet(packet_id, args.packets)
        tx_timestamp_ns = packet_id * period_ns
        lost = is_lost(rng, scenario)

        if lost:
            rx_timestamp_ns = ""
            one_way_delay_ns = ""
            jitter_ns = ""
        else:
            delay_ns = max(0, delay_for_scenario(rng, scenario, base_delay_ns, baseline_jitter_ns))
            rx_timestamp_ns = tx_timestamp_ns + delay_ns
            one_way_delay_ns = delay_ns
            if previous_received_delay is None:
                jitter_ns = 0
            else:
                jitter_ns = delay_ns - previous_received_delay
            previous_received_delay = delay_ns

        rows.append(
            {
                "packet_id": packet_id,
                "tx_timestamp_ns": tx_timestamp_ns,
                "rx_timestamp_ns": rx_timestamp_ns,
                "one_way_delay_ns": one_way_delay_ns,
                "jitter_ns": jitter_ns,
                "lost": int(lost),
                "burst_id": burst_id,
                "scenario": scenario,
            }
        )

    return rows


def write_csv(rows: list[dict[str, object]], output_path: Path) -> None:
    output_path.parent.mkdir(parents=True, exist_ok=True)
    fieldnames = [
        "packet_id",
        "tx_timestamp_ns",
        "rx_timestamp_ns",
        "one_way_delay_ns",
        "jitter_ns",
        "lost",
        "burst_id",
        "scenario",
    ]
    with output_path.open("w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(rows)


def main() -> int:
    args = parse_args()
    rows = generate_rows(args)
    output_path = Path(args.output)
    write_csv(rows, output_path)
    print(f"Synthetic SLA trace written to: {output_path}")
    print(f"Packets: {len(rows)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
