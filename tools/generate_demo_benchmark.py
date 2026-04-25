#!/usr/bin/env python3
"""Generate a deterministic demonstration benchmark dataset and SVG charts.

The generated data is synthetic but shaped like a real acceptance run:
- stable throughput around a 1G Ethernet service target;
- sub-ms one-way delay;
- bounded jitter;
- rare packet loss events.
"""

from __future__ import annotations

import csv
import math
import random
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DATA = ROOT / "results" / "demo-benchmark" / "metrics.csv"
SVG = ROOT / "docs" / "assets" / "generated_benchmark.svg"


def generate_rows() -> list[dict[str, float]]:
    random.seed(2544)
    rows = []
    for t in range(0, 61):
        delay_ms = 0.384 + 0.018 * math.sin(t / 7.0) + random.uniform(-0.006, 0.006)
        jitter_us = 24 + 10 * abs(math.sin(t / 5.5)) + random.uniform(0, 7)
        throughput_mbps = 941 + 8 * math.sin(t / 9.0) + random.uniform(-4, 4)
        loss_percent = 0.0
        if t in (17, 41):
            loss_percent = 0.02
        rows.append(
            {
                "time_s": t,
                "delay_ms": round(delay_ms, 6),
                "jitter_us": round(jitter_us, 3),
                "throughput_mbps": round(throughput_mbps, 3),
                "loss_percent": round(loss_percent, 4),
            }
        )
    return rows


def write_csv(rows: list[dict[str, float]]) -> None:
    DATA.parent.mkdir(parents=True, exist_ok=True)
    with DATA.open("w", newline="") as f:
        writer = csv.DictWriter(f, fieldnames=list(rows[0].keys()))
        writer.writeheader()
        writer.writerows(rows)


def polyline(rows: list[dict[str, float]], key: str, x0: int, y0: int, w: int, h: int) -> str:
    values = [float(r[key]) for r in rows]
    mn, mx = min(values), max(values)
    span = mx - mn or 1.0
    points = []
    for i, v in enumerate(values):
        x = x0 + int(w * i / (len(values) - 1))
        y = y0 + h - int(h * (v - mn) / span)
        points.append(f"{x},{y}")
    return " ".join(points)


def write_svg(rows: list[dict[str, float]]) -> None:
    SVG.parent.mkdir(parents=True, exist_ok=True)
    delay = polyline(rows, "delay_ms", 80, 150, 440, 120)
    jitter = polyline(rows, "jitter_us", 80, 370, 440, 120)
    thr = polyline(rows, "throughput_mbps", 650, 150, 440, 120)
    loss = polyline(rows, "loss_percent", 650, 370, 440, 120)
    SVG.write_text(
        f'''<svg xmlns="http://www.w3.org/2000/svg" width="1180" height="560" viewBox="0 0 1180 560">
<rect width="1180" height="560" rx="24" fill="#0b1220"/>
<text x="40" y="55" font-family="Arial" font-size="28" font-weight="700" fill="#e5e7eb">Generated benchmark from CSV</text>
<text x="40" y="82" font-family="Arial" font-size="15" fill="#94a3b8">Synthetic deterministic dataset for documentation, CI checks and report rendering experiments</text>
<g fill="#111827" stroke="#334155"><rect x="50" y="115" width="500" height="185" rx="16"/><rect x="620" y="115" width="500" height="185" rx="16"/><rect x="50" y="335" width="500" height="185" rx="16"/><rect x="620" y="335" width="500" height="185" rx="16"/></g>
<g font-family="Arial" fill="#e5e7eb" font-size="16" font-weight="700"><text x="75" y="142">One-way delay, ms</text><text x="645" y="142">Throughput, Mbit/s</text><text x="75" y="362">Jitter, us</text><text x="645" y="362">Loss, %</text></g>
<polyline points="{delay}" fill="none" stroke="#38bdf8" stroke-width="3"/>
<polyline points="{thr}" fill="none" stroke="#22c55e" stroke-width="3"/>
<polyline points="{jitter}" fill="none" stroke="#facc15" stroke-width="3"/>
<polyline points="{loss}" fill="none" stroke="#fb7185" stroke-width="3"/>
<g stroke="#334155"><line x1="80" y1="270" x2="520" y2="270"/><line x1="80" y1="490" x2="520" y2="490"/><line x1="650" y1="270" x2="1090" y2="270"/><line x1="650" y1="490" x2="1090" y2="490"/></g>
</svg>''',
        encoding="utf-8",
    )


def main() -> None:
    rows = generate_rows()
    write_csv(rows)
    write_svg(rows)
    print(f"Wrote {DATA}")
    print(f"Wrote {SVG}")


if __name__ == "__main__":
    main()
