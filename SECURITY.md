# Security Policy

## Supported Scope

This repository is maintained as a public demonstration and reference implementation. Security fixes should target the current `main` branch.

## Reporting A Vulnerability

Please do not open a public issue for a suspected vulnerability. Contact the repository owner privately and include:

- affected component;
- steps to reproduce;
- expected impact;
- suggested mitigation, if known.

## Secrets

Do not commit real credentials, `.env` files, database dumps with production data, logs, packet captures, or customer-specific configuration. Use `.env.example` for placeholders only.
