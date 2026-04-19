# Server SFP SLA

Server SFP SLA is a Linux service and PHP web console for configuring SFP-SLA modules, running SLA/RFC 2544/Y.1564 measurements, generating raw Ethernet traffic, and storing results in MySQL or MariaDB.

The repository is structured as a self-contained product demo: Go service, database schema, web console, deployment scripts, and operational documentation live in one place.

## Features

- Raw Ethernet traffic generation for SFP-SLA measurement scenarios.
- SLA real-time tests with delay, jitter, one-way delay, and packet-loss thresholds.
- RFC 2544-style throughput, latency, frame-loss, and burst tests.
- ITU-T Y.1564 service activation test workflow.
- PHP web console for modules, tests, live status, and result charts.
- MySQL/MariaDB persistence shared by the service and web console.
- CSRF protection, centralized web authentication, prepared statements for write endpoints, and hashed web-user passwords.

## Repository Layout

- `cmd/server-sfp-sla` - Go service entry point and measurement workflows.
- `internal/zsocket` - Linux packet socket/ring-buffer implementation used by traffic generators.
- `web/htdocs` - PHP web console served by Apache.
- `db/server_sfp_sla.sql` - MySQL/MariaDB schema and bootstrap data.
- `scripts` - installation, deployment, user management, cleanup, and validation scripts.
- `deploy/systemd` - systemd unit for the Go service.
- `docs/setup.md` - installation and operations guide.
- `docs/web.md` - web-console architecture and security notes.

## Requirements

- Linux host with root privileges for deployment and raw socket access.
- Go 1.22+ for building from source.
- MySQL or MariaDB.
- Apache with PHP and the `mysqli` and `snmp` extensions.
- `libpcap0.8-dev` and a network interface that is allowed to generate raw Ethernet traffic.

## Configuration

The service and the web console read database configuration from environment variables:

```bash
export SFP_SLA_DB_USER=sfp_user
export SFP_SLA_DB_PASSWORD='<change-me>'
export SFP_SLA_DB_NAME=server_sfp_sla
export SFP_SLA_DB_ADDR=127.0.0.1:3306
export SFP_SLA_TIMEZONE=UTC
```

If `SFP_SLA_DB_ADDR` is empty, the Go service uses the local MySQL socket. The PHP console can also use `SFP_SLA_DB_HOST` and `SFP_SLA_DB_PORT` when a socket is not desired. See `.env.example` for the full list.

Do not publish real credentials. The sample values are placeholders only.

## Build

```bash
sudo apt update
sudo apt install -y build-essential libpcap0.8-dev
go mod download
go build -o build/Server_SFP_SLA ./cmd/server-sfp-sla
```

## Install

```bash
export MYSQL_ROOT_PASSWORD='<mysql-root-password>'
export SFP_SLA_DB_PASSWORD='<application-db-password>'
sudo -E scripts/init.sh
sudo -E scripts/configServer.sh
sudo -E scripts/add_user.sh
```

`scripts/configServer.sh` imports `db/server_sfp_sla.sql`, deploys the Go service, and publishes `web/htdocs` to `/var/www/html`. `scripts/add_user.sh` creates the first web-console administrator with a `password_hash()` value.

## Web Security

- All `web/htdocs/action_*.php` files use the shared `web/htdocs/app.php` layer for authentication, CSRF checks, database access, and prepared statements.
- POST forms receive a hidden `csrf_token` field automatically.
- Web-user passwords are stored with PHP `password_hash()`.
- Public repository data does not include a default web-console password.

## Validation

```bash
make test
make lint
```

Equivalent manual checks:

```bash
go test ./...
go build -o /tmp/Server_SFP_SLA_check ./cmd/server-sfp-sla
bash -n scripts/*.sh
python scripts/check_web_links.py
find web/htdocs -name '*.php' -print0 | xargs -0 -n1 php -l
```

PHP and Go tooling must be installed locally for the corresponding checks.

## Public Demo Notes

- The schema contains no production secrets or default web login.
- Runtime credentials are supplied through environment variables or `/etc/default/server-sfp-sla`.
- Generated binaries, archives, local `.env` files, logs, editor settings, and extracted temporary artifacts are ignored.
- The web asset tree has been pruned to the files that are referenced by the PHP pages and CSS.
- GitHub Actions CI validates Go tests/builds, shell syntax, PHP syntax, and local web asset links.

