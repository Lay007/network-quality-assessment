#!/usr/bin/env python3
"""Check local Markdown links and image references.

The checker intentionally focuses on repository-local links because external
URLs may be unavailable in CI or depend on network policy.
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
LINK_RE = re.compile(r"!?\[[^\]]*\]\(([^)]+)\)")


def is_external(target: str) -> bool:
    return target.startswith(("http://", "https://", "mailto:", "#"))


def normalize(target: str) -> str:
    return target.split("#", 1)[0].strip()


def check_file(path: Path) -> list[str]:
    errors: list[str] = []
    text = path.read_text(encoding="utf-8")
    for match in LINK_RE.finditer(text):
        raw = match.group(1).strip()
        target = normalize(raw)
        if not target or is_external(raw):
            continue
        candidate = (path.parent / target).resolve()
        try:
            candidate.relative_to(ROOT)
        except ValueError:
            errors.append(f"{path.relative_to(ROOT)}: link escapes repo: {raw}")
            continue
        if not candidate.exists():
            errors.append(f"{path.relative_to(ROOT)}: broken link: {raw}")
    return errors


def main() -> int:
    errors: list[str] = []
    for md in ROOT.rglob("*.md"):
        if ".git" in md.parts:
            continue
        errors.extend(check_file(md))
    if errors:
        print("Broken Markdown links found:")
        for err in errors:
            print(f"- {err}")
        return 1
    print("All repository-local Markdown links are valid.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
