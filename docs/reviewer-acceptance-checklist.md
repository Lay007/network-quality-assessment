# Reviewer Acceptance Checklist

This checklist defines what a reviewer should verify before treating a network-quality measurement result as credible.

## Minimum credible result

| Area | Acceptance question | Evidence to look for |
|---|---|---|
| Trace schema | Are the packet timestamps and identifiers well defined? | `docs/sla-trace-schema.md`, CSV header, packet sequence fields |
| Reproducibility | Can the analysis be rerun from a clean checkout? | Synthetic trace generation and `tools/analyze_sla_trace.py` commands |
| Delay metric | Is one-way delay separated from host-side noise where possible? | timestamp source notes, hardware/software timestamp comparison |
| Jitter metric | Is jitter computed from a consistent delay series? | summary CSV, jitter histogram, report text |
| Packet loss | Are lost packets derived from sequence gaps rather than visual guessing? | packet-loss timeline, missing sequence accounting |
| Clock assumptions | Are PTP/GPS/local-clock limits documented? | clock synchronization notes and timing error budget |
| Report conclusion | Does the report state pass/fail/unknown with limitations? | generated `report.md` or filled SLA report |

## Review route

1. Open `README.md` and check the synthetic SLA demo commands.
2. Open `docs/sla-trace-schema.md` and verify that the fields match the generated CSV.
3. Run the synthetic demo locally.
4. Open `verification/reports/synthetic_sla_demo/report.md`.
5. Check the plots for delay trend, jitter distribution and packet-loss timeline.
6. Compare the final conclusion with the numeric summary.

## Red flags

A measurement should not be promoted as engineering evidence when:

- timestamps are present but the timestamp source is not described;
- one-way delay is reported without any clock assumption;
- packet loss is estimated visually instead of by sequence accounting;
- plots exist but the raw trace or schema is missing;
- the report has no limitations section;
- the result cannot be regenerated or traced to a specific input file.

## Recommended conclusion wording

Use clear language:

```text
The measurement is accepted for synthetic/demo validation because the trace schema, analysis command, generated plots and summary report are reproducible.
It is not yet accepted as hardware SLA evidence until the timestamp source, clock synchronization and physical test topology are documented.
```

This distinction is important: a good synthetic demo proves the reporting pipeline. It does not by itself prove the physical network path.
