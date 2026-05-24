# Release checklist

Use this checklist before publishing a public release of `network-quality-assessment`.

## 1. Local validation

```bash
go test ./...
go test -cover ./...
go vet ./...
test -z "$(gofmt -l cmd internal)"
bash -n scripts/*.sh
find web/htdocs -name '*.php' -print0 | xargs -0 -n1 php -l
python scripts/check_web_links.py
```

## 2. Demo and report artifacts

Generate or refresh the public demo artifacts:

```bash
python tools/generate_demo_benchmark.py
```

Check that demo outputs remain consistent with the README and report templates.

## 3. Measurement documentation review

- Timestamping model is clear.
- Software versus hardware timestamp assumptions are explicit.
- Units are stated for delay, jitter, loss and throughput.
- SLA thresholds are documented.
- Timing error budget is linked from the README.
- Demo scenario is reproducible without private traffic captures.

## 4. Release notes draft

Suggested release summary:

```text
v0.1.0 introduces a reproducible network measurement methodology repository with architecture diagrams, packet flow, FPGA/SFP timestamping concept, SLA report template, demo benchmark artifacts and credibility notes for software versus hardware timestamping.
```

## 5. Final checks

- CI is green.
- CodeQL is green or known issues are documented.
- README quick navigation points to current files.
- No private captures, credentials, topology details or production configs are committed.
- `CHANGELOG.md` has a release section.
