# Measurement quality gates

Use this checklist before merging changes that affect measurements, reports, dashboards, timestamp handling, or CI.

## Measurement credibility

- State which metric is affected: latency, jitter, packet loss, throughput, timestamp error, or report completeness.
- Document timestamp source and clock assumptions.
- Keep packet filters and field extraction rules visible.
- Explain whether a value comes from software timestamps, captured packets, logs, or generated summaries.
- Avoid mixing data from different environments in one before/after claim.

## Reproducibility

- Keep demo inputs small enough for CI or document where external data comes from.
- Provide exact command lines for generating reports.
- Record tool versions for packet capture, report generation, and graph generation.
- Store generated examples only when they are deterministic or clearly marked as illustrative.

## Validation

- Run Go tests when Go code changes.
- Run formatting and vet checks for Go changes.
- Run shell or PHP lint when those files change.
- Regenerate demo reports when report templates or plotting logic change.
- Check documentation links when README or docs are updated.

## Data safety

- Do not commit private captures, production addresses, credentials, or customer identifiers.
- Prefer anonymized or synthetic packet/report examples.
- Keep secrets out of examples and screenshots.

## Review checklist

- The changed metric is named.
- The validation command is documented.
- Before/after evidence is attached when behavior changes.
- The measurement environment is described well enough to reproduce the result.
