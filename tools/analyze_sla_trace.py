#!/usr/bin/env python3
"""Analyze a synthetic SLA trace and generate CSV/SVG/Markdown artifacts.

No third-party dependencies are required. SVG plots are generated directly so the
pipeline can run in CI without matplotlib.
"""

from __future__ import annotations

import argparse
import csv
import math
from collections import defaultdict
from pathlib import Path
from statistics import mean
from xml.sax.saxutils import escape


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Analyze a synthetic SLA packet trace")
    parser.add_argument("--input", required=True, help="Input synthetic_trace.csv")
    parser.add_argument("--output-dir", required=True, help="Output directory for reports")
    parser.add_argument("--delay-p95-threshold-us", type=float, default=450.0)
    parser.add_argument("--jitter-p95-threshold-us", type=float, default=50.0)
    parser.add_argument("--loss-threshold-pct", type=float, default=1.0)
    return parser.parse_args()


def read_trace(path: Path) -> list[dict[str, object]]:
    rows: list[dict[str, object]] = []
    with path.open("r", newline="", encoding="utf-8") as f:
        reader = csv.DictReader(f)
        for row in reader:
            parsed: dict[str, object] = {
                "packet_id": int(row["packet_id"]),
                "tx_timestamp_ns": int(row["tx_timestamp_ns"]),
                "lost": int(row["lost"]),
                "burst_id": int(row["burst_id"]),
                "scenario": row["scenario"],
            }
            for field in ("rx_timestamp_ns", "one_way_delay_ns", "jitter_ns"):
                parsed[field] = None if row[field] == "" else int(row[field])
            rows.append(parsed)
    return rows


def percentile(values: list[float], p: float) -> float:
    if not values:
        return float("nan")
    xs = sorted(values)
    if len(xs) == 1:
        return xs[0]
    rank = (len(xs) - 1) * p
    lo = math.floor(rank)
    hi = math.ceil(rank)
    if lo == hi:
        return xs[lo]
    return xs[lo] * (hi - rank) + xs[hi] * (rank - lo)


def scenario_summary(rows: list[dict[str, object]], args: argparse.Namespace) -> list[dict[str, object]]:
    groups: dict[str, list[dict[str, object]]] = defaultdict(list)
    for row in rows:
        groups[str(row["scenario"])].append(row)

    preferred_order = ["baseline", "jitter_burst", "loss_burst", "clock_offset", "mixed_impairments"]
    names = [name for name in preferred_order if name in groups]
    names.extend(sorted(set(groups) - set(names)))

    out: list[dict[str, object]] = []
    for name in names:
        sub = groups[name]
        received = [r for r in sub if int(r["lost"]) == 0]
        delays_us = [float(r["one_way_delay_ns"]) / 1000.0 for r in received if r["one_way_delay_ns"] is not None]
        jitter_us = [abs(float(r["jitter_ns"]) / 1000.0) for r in received if r["jitter_ns"] is not None]
        packets_total = len(sub)
        packets_lost = sum(int(r["lost"]) for r in sub)
        loss_pct = 100.0 * packets_lost / max(packets_total, 1)
        delay_mean_us = mean(delays_us) if delays_us else float("nan")
        delay_p95_us = percentile(delays_us, 0.95)
        jitter_p95_us = percentile(jitter_us, 0.95)

        sla_pass = (
            loss_pct <= args.loss_threshold_pct
            and delay_p95_us <= args.delay_p95_threshold_us
            and jitter_p95_us <= args.jitter_p95_threshold_us
        )
        root_cause_hint = root_cause(loss_pct, delay_p95_us, jitter_p95_us, args)

        out.append(
            {
                "scenario": name,
                "packets_total": packets_total,
                "packets_lost": packets_lost,
                "loss_pct": round(loss_pct, 4),
                "delay_mean_us": round(delay_mean_us, 3),
                "delay_p95_us": round(delay_p95_us, 3),
                "jitter_p95_us": round(jitter_p95_us, 3),
                "sla_pass": str(bool(sla_pass)).lower(),
                "root_cause_hint": root_cause_hint,
            }
        )
    return out


def root_cause(loss_pct: float, delay_p95_us: float, jitter_p95_us: float, args: argparse.Namespace) -> str:
    if loss_pct > args.loss_threshold_pct and jitter_p95_us > args.jitter_p95_threshold_us:
        return "loss and jitter impairment"
    if loss_pct > args.loss_threshold_pct:
        return "packet loss burst"
    if jitter_p95_us > args.jitter_p95_threshold_us:
        return "jitter burst or queueing"
    if delay_p95_us > args.delay_p95_threshold_us:
        return "delay bias or clock/path offset"
    return "within SLA thresholds"


def write_summary_csv(summary: list[dict[str, object]], path: Path) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    fieldnames = [
        "scenario",
        "packets_total",
        "packets_lost",
        "loss_pct",
        "delay_mean_us",
        "delay_p95_us",
        "jitter_p95_us",
        "sla_pass",
        "root_cause_hint",
    ]
    with path.open("w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(summary)


def scale(value: float, vmin: float, vmax: float, out_min: float, out_max: float) -> float:
    if vmax <= vmin:
        return (out_min + out_max) / 2.0
    t = (value - vmin) / (vmax - vmin)
    return out_min + t * (out_max - out_min)


def svg_header(width: int, height: int, title: str) -> list[str]:
    return [
        f'<svg width="{width}" height="{height}" viewBox="0 0 {width} {height}" xmlns="http://www.w3.org/2000/svg" role="img">',
        '<rect width="100%" height="100%" fill="#ffffff"/>',
        f'<text x="24" y="34" font-family="Arial, sans-serif" font-size="22" font-weight="700" fill="#102033">{escape(title)}</text>',
    ]


def write_delay_plot(rows: list[dict[str, object]], path: Path) -> None:
    received = [r for r in rows if int(r["lost"]) == 0]
    points = [(int(r["packet_id"]), float(r["one_way_delay_ns"]) / 1000.0) for r in received if r["one_way_delay_ns"] is not None]
    width, height = 1100, 520
    left, top, right, bottom = 70, 60, 30, 70
    xs = [p[0] for p in points]
    ys = [p[1] for p in points]
    x_min, x_max = min(xs), max(xs)
    y_min = max(0.0, min(ys) - 20.0)
    y_max = max(ys) + 20.0
    lines = svg_header(width, height, "Synthetic SLA Demo: one-way delay")
    plot_w = width - left - right
    plot_h = height - top - bottom
    lines.append(f'<rect x="{left}" y="{top}" width="{plot_w}" height="{plot_h}" fill="#f7fbff" stroke="#c8d7e6"/>')

    poly = []
    for packet_id, delay_us in points:
        x = scale(packet_id, x_min, x_max, left, left + plot_w)
        y = scale(delay_us, y_min, y_max, top + plot_h, top)
        poly.append(f"{x:.2f},{y:.2f}")
    lines.append(f'<polyline points="{" ".join(poly)}" fill="none" stroke="#2364aa" stroke-width="1.6"/>')
    add_axes(lines, left, top, plot_w, plot_h, "packet_id", "delay, us")
    lines.append("</svg>")
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_jitter_plot(rows: list[dict[str, object]], path: Path) -> None:
    values = [abs(float(r["jitter_ns"]) / 1000.0) for r in rows if int(r["lost"]) == 0 and r["jitter_ns"] is not None]
    width, height = 900, 520
    left, top, right, bottom = 70, 60, 30, 70
    bins = 30
    vmax = max(values) if values else 1.0
    counts = [0 for _ in range(bins)]
    for value in values:
        idx = min(bins - 1, int((value / max(vmax, 1e-9)) * bins))
        counts[idx] += 1
    max_count = max(counts) if counts else 1
    plot_w = width - left - right
    plot_h = height - top - bottom
    bar_w = plot_w / bins
    lines = svg_header(width, height, "Synthetic SLA Demo: absolute jitter histogram")
    lines.append(f'<rect x="{left}" y="{top}" width="{plot_w}" height="{plot_h}" fill="#f7fbff" stroke="#c8d7e6"/>')
    for i, count in enumerate(counts):
        h = 0 if max_count == 0 else plot_h * count / max_count
        x = left + i * bar_w
        y = top + plot_h - h
        lines.append(f'<rect x="{x:.2f}" y="{y:.2f}" width="{max(bar_w - 2, 1):.2f}" height="{h:.2f}" fill="#2364aa" opacity="0.82"/>')
    add_axes(lines, left, top, plot_w, plot_h, "absolute jitter, us", "count")
    lines.append("</svg>")
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_loss_plot(rows: list[dict[str, object]], path: Path) -> None:
    width, height = 1100, 300
    left, top, right, bottom = 70, 60, 30, 60
    plot_w = width - left - right
    plot_h = height - top - bottom
    packet_ids = [int(r["packet_id"]) for r in rows]
    x_min, x_max = min(packet_ids), max(packet_ids)
    lines = svg_header(width, height, "Synthetic SLA Demo: packet loss timeline")
    lines.append(f'<rect x="{left}" y="{top}" width="{plot_w}" height="{plot_h}" fill="#f7fbff" stroke="#c8d7e6"/>')
    for row in rows:
        if int(row["lost"]) == 0:
            continue
        x = scale(int(row["packet_id"]), x_min, x_max, left, left + plot_w)
        lines.append(f'<line x1="{x:.2f}" y1="{top}" x2="{x:.2f}" y2="{top + plot_h}" stroke="#cc2936" stroke-width="1.8"/>')
    add_axes(lines, left, top, plot_w, plot_h, "packet_id", "loss events")
    lines.append("</svg>")
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def add_axes(lines: list[str], left: int, top: int, plot_w: int, plot_h: int, xlabel: str, ylabel: str) -> None:
    lines.append(f'<line x1="{left}" y1="{top + plot_h}" x2="{left + plot_w}" y2="{top + plot_h}" stroke="#102033"/>')
    lines.append(f'<line x1="{left}" y1="{top}" x2="{left}" y2="{top + plot_h}" stroke="#102033"/>')
    lines.append(f'<text x="{left + plot_w / 2:.1f}" y="{top + plot_h + 45}" font-family="Arial, sans-serif" font-size="15" text-anchor="middle" fill="#102033">{escape(xlabel)}</text>')
    lines.append(f'<text x="22" y="{top + plot_h / 2:.1f}" font-family="Arial, sans-serif" font-size="15" text-anchor="middle" transform="rotate(-90 22 {top + plot_h / 2:.1f})" fill="#102033">{escape(ylabel)}</text>')


def write_report(summary: list[dict[str, object]], output_dir: Path, input_path: Path, args: argparse.Namespace) -> None:
    lines = [
        "# Synthetic SLA Demo Report",
        "",
        "This report is generated from a deterministic synthetic packet trace. It demonstrates the reporting flow without FPGA/SFP timestamping hardware.",
        "",
        "## Input",
        "",
        f"- Trace: `{input_path}`",
        f"- Delay p95 threshold: `{args.delay_p95_threshold_us}` us",
        f"- Jitter p95 threshold: `{args.jitter_p95_threshold_us}` us",
        f"- Loss threshold: `{args.loss_threshold_pct}` %",
        "",
        "## SLA summary",
        "",
        "| Scenario | Packets | Lost | Loss, % | Delay mean, us | Delay p95, us | Jitter p95, us | SLA | Root cause hint |",
        "|---|---:|---:|---:|---:|---:|---:|---|---|",
    ]
    for row in summary:
        sla = "PASS" if row["sla_pass"] == "true" else "FAIL"
        lines.append(
            f"| `{row['scenario']}` | {row['packets_total']} | {row['packets_lost']} | {row['loss_pct']} | "
            f"{row['delay_mean_us']} | {row['delay_p95_us']} | {row['jitter_p95_us']} | {sla} | {row['root_cause_hint']} |"
        )
    lines.extend(
        [
            "",
            "## Generated artifacts",
            "",
            "- `synthetic_trace.csv`",
            "- `sla_summary.csv`",
            "- `one_way_delay_timeseries.svg`",
            "- `jitter_histogram.svg`",
            "- `packet_loss_timeline.svg`",
            "",
            "## Interpretation",
            "",
            "The synthetic trace is intended for CI, documentation and reviewer onboarding. It is not a substitute for real hardware timestamping. Real FPGA/SFP measurements should use the same report structure but must state clocking, timestamp insertion point, capture path and calibration assumptions.",
        ]
    )
    (output_dir / "report.md").write_text("\n".join(lines) + "\n", encoding="utf-8")


def main() -> int:
    args = parse_args()
    input_path = Path(args.input)
    output_dir = Path(args.output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)

    rows = read_trace(input_path)
    summary = scenario_summary(rows, args)
    write_summary_csv(summary, output_dir / "sla_summary.csv")
    write_delay_plot(rows, output_dir / "one_way_delay_timeseries.svg")
    write_jitter_plot(rows, output_dir / "jitter_histogram.svg")
    write_loss_plot(rows, output_dir / "packet_loss_timeline.svg")
    write_report(summary, output_dir, input_path, args)

    print(f"SLA summary written to: {output_dir / 'sla_summary.csv'}")
    print(f"Report written to: {output_dir / 'report.md'}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
