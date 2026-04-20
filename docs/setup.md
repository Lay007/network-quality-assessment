# Setup Guide

This guide describes how to install Server SFP SLA on a Linux host. The commands assume root privileges because the service uses raw sockets and is installed as a systemd service.

## Installation Flow

```bash
export MYSQL_ROOT_PASSWORD='<mysql-root-password>'
export SFP_SLA_DB_USER=sfp_user
export SFP_SLA_DB_PASSWORD='<application-db-password>'
export SFP_SLA_DB_NAME=server_sfp_sla
export SFP_SLA_DB_ADDR=''
export SFP_SLA_TIMEZONE=UTC

sudo -E scripts/init.sh
sudo -E scripts/configServer.sh
sudo -E scripts/add_user.sh
```

## What The Scripts Do

- `scripts/init.sh` installs system dependencies: Go build prerequisites, Apache, MySQL/MariaDB-compatible packages, PHP extensions, `net-tools`, and `pv`.
- `scripts/configServer.sh` creates the application database and database user, imports `db/server_sfp_sla.sql`, deploys the Go service, and publishes the web console to `/var/www/html`.
- `scripts/deploy.sh` builds the Go binary when needed, installs it to `/usr/local/bin/Server_SFP_SLA`, installs the systemd unit, creates `/etc/default/server-sfp-sla` when missing, and starts the service.
- `scripts/add_user.sh` creates a web-console user with a hashed password.
- `scripts/clean.sh` removes the installed database, database user, service files, web files, and packages installed by `scripts/init.sh`.
- `scripts/develop_config.sh` installs the build dependencies needed for local development on Linux.

## Required Files

- `db/server_sfp_sla.sql` - database schema and bootstrap configuration.
- `build/Server_SFP_SLA` - optional prebuilt service binary; if missing, `scripts/deploy.sh` builds it from source.
- `deploy/systemd/Server_SFP_SLA.service` - systemd unit.
- `web/htdocs` - Apache document root content for the web console.

## Configuration

Use environment variables for credentials and database parameters:

```bash
export MYSQL_ROOT_PASSWORD='<mysql-root-password>'
export SFP_SLA_DB_USER=sfp_user
export SFP_SLA_DB_PASSWORD='<application-db-password>'
export SFP_SLA_DB_NAME=server_sfp_sla
export SFP_SLA_DB_ADDR=127.0.0.1:3306
export SFP_SLA_TIMEZONE=UTC
```

`SFP_SLA_DB_ADDR` should be empty for a local MySQL socket connection. Use `host:port` for TCP.

For systemd deployments, move production values into `/etc/default/server-sfp-sla` and restrict that file to root-readable permissions.

## Post-Install Checks

```bash
systemctl status Server_SFP_SLA
systemctl status apache2
systemctl status mysql
```

Open the web console through the Apache server address. After logging in, select the network interface in **Global settings** before starting traffic-generating tests.

## Uninstall

```bash
sudo -E scripts/clean.sh
```

Review `scripts/clean.sh` before running it on a shared host. It removes application database objects and packages installed by the setup script.
