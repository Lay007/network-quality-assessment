# Web Console

`web/htdocs` contains the PHP web console. It shares the MySQL/MariaDB database with the Go service:

- The web console creates modules, global settings, and test jobs in `modules_sfp_sla`, `global_config`, and `test_*` tables.
- The Go service reads pending jobs, performs measurements, and writes results.
- The web console displays status, result tables, and charts.

## Shared Layer

- `web/htdocs/app.php` centralizes sessions, authentication, CSRF protection, database connection, escaping, and query helpers.
- `web/htdocs/db.php` is a compatibility shim that exposes `$db` and `$link` for older pages while still using `app.php`.
- `web/htdocs/start.php` is the public login form.
- `web/htdocs/notepad.php` handles login and uses `password_verify()`.

## Security

- `action_*.php` write endpoints require POST, authentication, CSRF, and prepared statements.
- Direct `$_POST` and `$_GET` reads in action endpoints are routed through `app_post_*()` and `app_get_*()` helpers.
- `jsonl.php` validates numeric parameters and uses prepared statements.
- Web-user passwords are stored as `password_hash()` values.
- Public schema data does not include a default administrator password.

## Asset Policy

The original admin template was reduced to referenced assets only. Unused skins, images, legacy IE polyfills, demo pages, unused plugins, and unreachable PHP/JSON endpoints were removed.

## Link Check

```bash
python scripts/check_web_links.py
```

The checker scans local `src`, `href`, `action`, and CSS `url(...)` references in `web/htdocs` and fails on missing files.

## Maintenance Notes

- The larger view files, such as `tests_view_all.php` and `module_view_all.php`, are still procedural PHP. The write endpoints are secured, but a future UI rewrite should move repeated layout blocks into shared includes.
- PHP lint requires the `php` CLI: `find web/htdocs -name '*.php' -print0 | xargs -0 -n1 php -l`.
