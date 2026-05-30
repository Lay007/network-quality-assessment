# Report review checklist

Use this checklist when changing generated reports, dashboards, screenshots, or demo artifacts.

## Report content

- The report states the measurement period or input dataset.
- Units are visible for every plotted or tabulated metric.
- Thresholds and SLA-style limits are explained.
- Missing data and invalid samples are handled explicitly.
- The report distinguishes measured values from derived or illustrative values.

## Graphs and tables

- Axes are readable in the README and rendered documentation.
- Legends do not overlap the plotted data.
- Tables use stable column names.
- Generated images can be traced back to source data or generation scripts.

## Evidence

- Include a small example report or screenshot when presentation changes.
- Include before/after evidence when calculations change.
- Mention expected differences across platforms or capture methods.

## Safety

- Remove private IP addresses, hostnames, usernames, and organization names from examples.
- Use synthetic or anonymized captures for public artifacts.
- Do not publish production packet captures.
