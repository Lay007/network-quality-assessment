<?php

if (session_status() === PHP_SESSION_NONE) {
    session_start();
}

date_default_timezone_set(getenv('SFP_SLA_TIMEZONE') ?: 'UTC');

function app_db(): mysqli
{
    static $db = null;
    if ($db instanceof mysqli) {
        return $db;
    }

    $addr = getenv('SFP_SLA_DB_ADDR') ?: '';
    $host = getenv('SFP_SLA_DB_HOST') ?: 'localhost';
    $port = getenv('SFP_SLA_DB_PORT') ?: null;
    if ($addr !== '') {
        $parts = explode(':', $addr, 2);
        $host = $parts[0];
        $port = $parts[1] ?? $port;
    }

    $database = getenv('SFP_SLA_DB_NAME') ?: 'server_sfp_sla';
    $user = getenv('SFP_SLA_DB_USER') ?: 'sfp_user';
    $password = getenv('SFP_SLA_DB_PASSWORD') ?: '';

    $db = mysqli_connect($host, $user, $password, $database, $port ? (int)$port : null);
    if (!$db) {
        http_response_code(500);
        exit('WRONG CONNECTION');
    }

    mysqli_set_charset($db, 'utf8');
    return $db;
}

function app_current_file(): string
{
    return basename($_SERVER['SCRIPT_NAME'] ?? '');
}

function app_is_public_page(): bool
{
    return in_array(app_current_file(), ['index.php', 'start.php', 'notepad.php'], true);
}

function app_is_authenticated(): bool
{
    return isset($_SESSION['user_in']) && $_SESSION['user_in'] !== '';
}

function app_redirect(string $url): void
{
    header('Location: ' . $url);
    exit;
}

function app_csrf_token(): string
{
    if (empty($_SESSION['csrf_token'])) {
        $_SESSION['csrf_token'] = bin2hex(random_bytes(32));
    }

    return $_SESSION['csrf_token'];
}

function app_validate_csrf(): void
{
    $token = $_POST['csrf_token'] ?? '';
    if (!is_string($token) || !hash_equals(app_csrf_token(), $token)) {
        http_response_code(403);
        exit('CSRF validation failed');
    }
}

function app_require_auth(): void
{
    if (!app_is_authenticated()) {
        app_redirect('start.php');
    }
}

function app_require_post(): void
{
    if (($_SERVER['REQUEST_METHOD'] ?? '') !== 'POST') {
        http_response_code(405);
        exit('Method not allowed');
    }

    app_validate_csrf();
    app_require_auth();
}

function app_post_string(string $name, string $default = ''): string
{
    $value = $_POST[$name] ?? $default;
    return is_string($value) ? trim($value) : $default;
}

function app_post_int(string $name, int $default = 0): int
{
    $value = $_POST[$name] ?? null;
    if ($value === null || $value === '') {
        return $default;
    }

    return (int)$value;
}

function app_post_float(string $name, float $default = 0.0): float
{
    $value = $_POST[$name] ?? null;
    if ($value === null || $value === '') {
        return $default;
    }

    return (float)$value;
}

function app_post_bool(string $name): int
{
    if (!isset($_POST[$name])) {
        return 0;
    }

    $value = $_POST[$name];
    if ($value === '' || $value === 'on') {
        return 1;
    }

    return in_array(strtolower((string)$value), ['1', 'true', 'yes'], true) ? 1 : 0;
}

function app_get_string(string $name, string $default = ''): string
{
    $value = $_GET[$name] ?? $default;
    return is_string($value) ? trim($value) : $default;
}

function app_get_int(string $name, int $default = 0): int
{
    $value = $_GET[$name] ?? null;
    if ($value === null || $value === '') {
        return $default;
    }

    return (int)$value;
}

function app_h($value): string
{
    return htmlspecialchars((string)$value, ENT_QUOTES, 'UTF-8');
}

function app_module_label(int $id): string
{
    if ($id <= 0) {
        return '';
    }

    $stmt = app_stmt(
        'SELECT `name`, `address_ip` FROM `modules_sfp_sla` WHERE `id` = ?',
        'i',
        [$id]
    );
    $result = mysqli_stmt_get_result($stmt);
    $row = $result ? mysqli_fetch_assoc($result) : null;
    if (!$row) {
        return '';
    }

    return app_h($row['name']) . ' : IP= ' . app_h($row['address_ip']);
}

function app_module_row(int $id): ?array
{
    if ($id <= 0) {
        return null;
    }

    $stmt = app_stmt(
        'SELECT * FROM `modules_sfp_sla` WHERE `id` = ?',
        'i',
        [$id]
    );
    $result = mysqli_stmt_get_result($stmt);
    $row = $result ? mysqli_fetch_assoc($result) : null;

    return $row ?: null;
}

function app_bind_refs(array &$values): array
{
    $refs = [];
    foreach ($values as $key => &$value) {
        $refs[$key] = &$value;
    }

    return $refs;
}

function app_stmt(string $sql, string $types = '', array $params = []): mysqli_stmt
{
    $stmt = mysqli_prepare(app_db(), $sql);
    if (!$stmt) {
        http_response_code(500);
        exit('SQL prepare failed');
    }

    if ($types !== '') {
        $refs = app_bind_refs($params);
        mysqli_stmt_bind_param($stmt, $types, ...$refs);
    }

    if (!mysqli_stmt_execute($stmt)) {
        http_response_code(500);
        exit('SQL execute failed');
    }

    return $stmt;
}

function app_scalar(string $sql, string $types = '', array $params = [])
{
    $stmt = app_stmt($sql, $types, $params);
    $result = mysqli_stmt_get_result($stmt);
    if (!$result) {
        return null;
    }

    $row = mysqli_fetch_row($result);
    return $row[0] ?? null;
}

function app_inject_csrf(string $buffer): string
{
    if (stripos($buffer, '<form') === false) {
        return $buffer;
    }

    $field = '<input type="hidden" name="csrf_token" value="' . htmlspecialchars(app_csrf_token(), ENT_QUOTES, 'UTF-8') . '">';

    return preg_replace_callback(
        '/<form\b[^>]*>/i',
        static function (array $matches) use ($field): string {
            $tag = $matches[0];
            if (!preg_match('/\bmethod\s*=\s*([\'"]?)post\1/i', $tag)) {
                return $tag;
            }

            return $tag . $field;
        },
        $buffer
    );
}

if (($_SERVER['REQUEST_METHOD'] ?? '') === 'POST') {
    app_validate_csrf();
}

if (!app_is_public_page()) {
    app_require_auth();
}

if (!headers_sent()) {
    ob_start('app_inject_csrf');
}
