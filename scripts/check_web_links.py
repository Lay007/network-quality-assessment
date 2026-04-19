#!/usr/bin/env python3
"""Validate local web asset references.

The checker intentionally stays simple: it scans PHP, HTML, CSS, and JS files for
static local references and reports paths that do not exist under web/htdocs.
"""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
WEB_ROOT = ROOT / "web" / "htdocs"
SCAN_EXTENSIONS = {".php", ".html", ".css", ".js"}
SKIP_DIRECTORIES = {"code", "vendor"}
ATTRIBUTE_REF = re.compile(r"""(?:href|src|action)\s*=\s*\\?["']([^"']+)""", re.IGNORECASE)
CSS_URL_REF = re.compile(r"""url\(\s*["']?([^"')]+)""", re.IGNORECASE)
SKIP_PREFIXES = ("#", "http:", "https:", "//", "data:", "javascript:", "mailto:")


def normalize(raw: str) -> str | None:
    value = raw.strip().replace('\\"', "").replace("\\'", "")
    if not value or value.startswith(SKIP_PREFIXES):
        return None
    value = value.split("?", 1)[0].split("#", 1)[0]
    return value or None


def target_for(source: Path, value: str) -> Path | None:
    target = WEB_ROOT / value.lstrip("/") if value.startswith("/") else source.parent / value
    target = target.resolve()
    try:
        target.relative_to(WEB_ROOT.resolve())
    except ValueError:
        return None
    return target


def main() -> int:
    missing: list[tuple[Path, str, Path]] = []
    ref_count = 0

    for source in sorted(WEB_ROOT.rglob("*")):
        if not source.is_file() or source.suffix.lower() not in SCAN_EXTENSIONS:
            continue
        if any(part in SKIP_DIRECTORIES for part in source.relative_to(WEB_ROOT).parts[:-1]):
            continue

        text = source.read_text(encoding="utf-8", errors="ignore")
        refs = [match.group(1) for match in ATTRIBUTE_REF.finditer(text)]
        if source.suffix.lower() == ".css":
            refs.extend(match.group(1) for match in CSS_URL_REF.finditer(text))

        for raw in refs:
            value = normalize(raw)
            if value is None:
                continue
            target = target_for(source, value)
            if target is None:
                continue
            ref_count += 1
            if not target.exists():
                missing.append((source.relative_to(WEB_ROOT), value, target.relative_to(WEB_ROOT)))

    print(f"ref_count {ref_count}")
    print(f"missing_count {len(missing)}")
    for source, value, target in missing:
        print(f"missing {source}: {value} -> {target}")

    return 1 if missing else 0


if __name__ == "__main__":
    sys.exit(main())
