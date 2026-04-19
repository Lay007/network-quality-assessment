<?php
require_once 'db.php';
header('Content-type: application/json');

$testId = max(0, app_get_int('id'));
$limit = min(10000, max(1, app_get_int('n', 600)));
$all = [];

$stmt = app_stmt(
    'SELECT `datatime`, `load_to_lazer`, `load_to_com` FROM `modules_sfp_sla_load_rez` WHERE `module_id` = ? AND datatime >= date_sub(now(), INTERVAL 5 MINUTE) ORDER BY `datatime` DESC LIMIT ?',
    'ii',
    [$testId, $limit]
);
$result = mysqli_stmt_get_result($stmt);
while ($row = mysqli_fetch_row($result)) {
    $all[] = [(strtotime($row[0])), (float)$row[1], (float)$row[2]];
}

echo json_encode($all);
