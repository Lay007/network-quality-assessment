#!/usr/bin/env python3
"""Run two built loopback agents and record software evidence without hardware."""
from __future__ import annotations

import argparse
import contextlib
import json
import os
from pathlib import Path
import platform
import re
import secrets
import statistics
import subprocess
import time


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--bin-dir", type=Path, default=Path("bin"))
    parser.add_argument("--output", type=Path, default=Path("bin/local-bench"))
    parser.add_argument("--runs", type=int, default=3)
    parser.add_argument("--packets", type=int, default=100)
    args = parser.parse_args()
    if not 1 <= args.runs <= 10 or not 1 <= args.packets <= 200:
        parser.error("runs must be 1..10 and packets 1..200")
    suffix = ".exe" if os.name == "nt" else ""
    agent = (args.bin_dir / ("measurement-agent" + suffix)).resolve()
    controller = (args.bin_dir / ("measurement-bench" + suffix)).resolve()
    if not agent.is_file() or not controller.is_file():
        parser.error("build measurement-agent and measurement-bench in --bin-dir first")
    args.output.mkdir(parents=True, exist_ok=True)
    env = dict(os.environ, NQA_BENCH_TOKEN=secrets.token_hex(24))
    evidence: dict[str, object] = {
        "schema_version": 1,
        "evidence_kind": "loopback-software-round-trip",
        "platform": {"os": platform.system(), "kernel": platform.release(), "architecture": platform.machine()},
        "method": "Two local agents exchange authenticated 64-byte UDP requests/replies, one outstanding per session, >=10ms pause after each completed probe. HTTP coordinates lifecycle, not synchronized clocks or launches.",
        "limitations": "Application RTT includes scheduling, queues and echo handling. Not one-way delay, hardware precision, line rate or a two-machine validation. RSS is the maximum individual child RSS, not aggregate memory.",
        "runs": [],
    }
    try:
        import resource
        before = resource.getrusage(resource.RUSAGE_CHILDREN)
    except ImportError:
        resource = None
        before = None
    processes: list[subprocess.Popen] = []
    start = time.perf_counter()
    graceful = True
    failed = False
    with contextlib.ExitStack() as stack:
        try:
            bases: list[str] = []
            for name in ("local-a", "local-b"):
                log_path = args.output / f"{name}.log"
                log = stack.enter_context(log_path.open("w", encoding="utf-8"))
                proc = subprocess.Popen([str(agent), "serve", "--id", name, "--listen", "127.0.0.1:0", "--probe", "127.0.0.1:0"], env=env, stdout=log, stderr=log)
                processes.append(proc)
                deadline = time.monotonic() + 10
                while True:
                    text = log_path.read_text(encoding="utf-8")
                    match = re.search(r"control=(http://127\.0\.0\.1:\d+)", text)
                    if match:
                        bases.append(match.group(1))
                        break
                    if proc.poll() is not None or time.monotonic() > deadline:
                        raise RuntimeError(f"{name} failed to start; inspect {log_path}")
                    time.sleep(0.02)
            for number in range(args.runs):
                completed = subprocess.run([str(controller), "--agents", ",".join(bases), "--packets", str(args.packets), "--output", str(args.output / "packages")], env=env, capture_output=True, text=True, timeout=15, check=True)
                package = Path(completed.stdout.strip())
                report = json.loads((package / "aggregate.json").read_text(encoding="utf-8"))
                samples = [sample["round_trip_ns"] for session in report["sessions"] for sample in session["samples"] if sample.get("round_trip_ns") is not None]
                if report["status"] != "completed" or any(session["result"]["timestamp_provenance"]["kind"] != "software" for session in report["sessions"]):
                    raise RuntimeError("unexpected report status or provenance")
                percentiles = statistics.quantiles(samples, n=100, method="inclusive") if len(samples) > 1 else samples * 99
                evidence["runs"].append({
                    "run": number + 1, "bench_session_id": report["bench_session_id"],
                    "packets_attempted": report["packets_attempted"], "packets_received": report["packets_received"], "packets_lost": report["packets_lost"],
                    "rtt_ns": {"min": min(samples), "mean": statistics.mean(samples), "p50": statistics.median(samples), "p95": percentiles[94], "p99": percentiles[98], "max": max(samples)} if samples else None,
                    "software_versions": sorted({agent["software_version"] for agent in report["agents"]}),
                })
        except (OSError, ValueError, RuntimeError, subprocess.SubprocessError) as exc:
            evidence["error"] = str(exc)
            failed = True
        finally:
            for proc in processes:
                proc.terminate()
            for proc in processes:
                try:
                    proc.wait(timeout=3)
                    graceful = graceful and proc.returncode == 0
                except subprocess.TimeoutExpired:
                    proc.kill()
                    proc.wait()
                    graceful = False
    evidence["elapsed_seconds"] = time.perf_counter() - start
    evidence["agents_exited_cleanly"] = graceful
    if resource is not None and before is not None:
        usage = resource.getrusage(resource.RUSAGE_CHILDREN)
        evidence["children_cpu_seconds"] = {"user": usage.ru_utime - before.ru_utime, "system": usage.ru_stime - before.ru_stime}
        evidence["max_child_rss_kib"] = usage.ru_maxrss if platform.system() == "Linux" else None
    else:
        evidence["children_cpu_seconds"] = None
        evidence["max_child_rss_kib"] = None
    path = args.output / "validation.json"
    path.write_text(json.dumps(evidence, indent=2) + "\n", encoding="utf-8")
    print(path)
    return int(failed or not graceful)


if __name__ == "__main__":
    raise SystemExit(main())
