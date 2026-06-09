#!/usr/bin/env python3
"""Lightweight tests for the synthetic SLA analyzer.

The tests intentionally use only the Python standard library so they can run in
GitHub Actions without additional dependencies.
"""

from __future__ import annotations

import argparse
import importlib.util
from pathlib import Path
from tempfile import TemporaryDirectory


ROOT = Path(__file__).resolve().parents[1]
ANALYZER = ROOT / "tools" / "analyze_sla_trace.py"


def load_analyzer():
    spec = importlib.util.spec_from_file_location("analyze_sla_trace", ANALYZER)
    if spec is None or spec.loader is None:
        raise RuntimeError(f"Cannot load analyzer from {ANALYZER}")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def make_args(delay: float = 450.0, jitter: float = 50.0, loss: float = 1.0) -> argparse.Namespace:
    return argparse.Namespace(
        delay_p95_threshold_us=delay,
        jitter_p95_threshold_us=jitter,
        loss_threshold_pct=loss,
    )


def test_percentile(analyzer) -> None:
    assert analyzer.percentile([], 0.95) != analyzer.percentile([], 0.95)  # NaN
    assert analyzer.percentile([10.0], 0.95) == 10.0
    assert analyzer.percentile([0.0, 10.0, 20.0, 30.0], 0.5) == 15.0
    assert analyzer.percentile([1.0, 2.0, 3.0, 4.0, 5.0], 0.95) == 4.8


def test_root_cause(analyzer) -> None:
    args = make_args()
    assert analyzer.root_cause(0.0, 320.0, 3.0, args) == "within SLA thresholds"
    assert analyzer.root_cause(0.0, 500.0, 3.0, args) == "delay bias or clock/path offset"
    assert analyzer.root_cause(0.0, 320.0, 80.0, args) == "jitter burst or queueing"
    assert analyzer.root_cause(5.0, 320.0, 3.0, args) == "packet loss burst"
    assert analyzer.root_cause(5.0, 320.0, 80.0, args) == "loss and jitter impairment"


def test_scenario_summary(analyzer) -> None:
    args = make_args(delay=450.0, jitter=50.0, loss=1.0)
    rows = []
    for packet_id in range(10):
        delay_ns = 320_000 + packet_id * 100
        rows.append(
            {
                "packet_id": packet_id,
                "tx_timestamp_ns": packet_id * 1_000_000,
                "rx_timestamp_ns": packet_id * 1_000_000 + delay_ns,
                "one_way_delay_ns": delay_ns,
                "jitter_ns": 0 if packet_id == 0 else 100,
                "lost": 0,
                "burst_id": 0,
                "scenario": "baseline",
            }
        )

    summary = analyzer.scenario_summary(rows, args)
    assert len(summary) == 1
    row = summary[0]
    assert row["scenario"] == "baseline"
    assert row["packets_total"] == 10
    assert row["packets_lost"] == 0
    assert row["sla_pass"] == "true"
    assert row["root_cause_hint"] == "within SLA thresholds"


def test_plot_rendering_with_empty_inputs(analyzer) -> None:
    with TemporaryDirectory() as tmp:
        out = Path(tmp)

        analyzer.write_delay_plot([], out / "delay-empty.svg")
        analyzer.write_loss_plot([], out / "loss-empty.svg")
        analyzer.write_delay_plot(
            [
                {
                    "packet_id": 1,
                    "tx_timestamp_ns": 0,
                    "rx_timestamp_ns": None,
                    "one_way_delay_ns": None,
                    "jitter_ns": None,
                    "lost": 1,
                    "burst_id": 0,
                    "scenario": "all_lost",
                }
            ],
            out / "delay-all-lost.svg",
        )

        assert "No packets available" in (out / "loss-empty.svg").read_text(encoding="utf-8")
        assert "No received packets available" in (out / "delay-empty.svg").read_text(encoding="utf-8")
        assert "No received packets available" in (out / "delay-all-lost.svg").read_text(encoding="utf-8")


def main() -> int:
    analyzer = load_analyzer()
    test_percentile(analyzer)
    test_root_cause(analyzer)
    test_scenario_summary(analyzer)
    test_plot_rendering_with_empty_inputs(analyzer)
    print("PASS: synthetic SLA analyzer unit tests")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
