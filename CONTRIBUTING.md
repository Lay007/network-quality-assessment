# Contributing

Thank you for improving `network-quality-assessment`.

This repository is an engineering project for network measurement methodology, timestamping, SLA-oriented metrics and reproducible reports. Contributions should improve measurement credibility, reproducibility, safety or operational clarity.

## Development principles

- Keep changes small and reviewable.
- Document measurement assumptions and units.
- Prefer reproducible synthetic datasets for demos and CI.
- Do not commit real credentials, production configs, private captures or customer data.
- Keep generated graphs and reports traceable to scripts and manifests.
- Update documentation when packet format, metrics, thresholds or report semantics change.

## Local validation

Recommended baseline:

```bash
go test ./...
go test -cover ./...
go vet ./...
test -z "$(gofmt -l cmd internal)"
bash -n scripts/*.sh
find web/htdocs -name '*.php' -print0 | xargs -0 -n1 php -l
python scripts/check_web_links.py
```

If available, use the repository `Makefile` target that matches the change.

## Measurement change checklist

For measurement logic, packet format, timestamping, reports or dashboards, include or update:

- metric definition;
- units and thresholds;
- synthetic demo data or manifest;
- expected report output;
- graph or dashboard artifact;
- error-budget note, if relevant.

## Documentation expectations

Documentation should answer:

- what is measured;
- where timestamps are taken;
- what noise sources exist;
- what assumptions are made;
- how the result can be reproduced;
- how software timestamps differ from hardware/FPGA timestamps.

## Pull request checklist

- [ ] Go tests pass.
- [ ] Go vet passes.
- [ ] Formatting is checked.
- [ ] Shell/PHP lint passes, if relevant.
- [ ] Web links are checked, if docs changed.
- [ ] Demo data or report artifacts are reproducible.
- [ ] No private data or credentials are committed.
