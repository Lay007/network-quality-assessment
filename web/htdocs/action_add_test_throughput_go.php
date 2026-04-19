<?php
require_once 'db.php';

app_require_post();

$testType = app_post_int('test_type');
$moduleFirst = app_post_int('module_first');
$moduleSecond = app_post_int('module_second');
$clock = app_post_int('clock');
$count = app_post_int('count');
$maxLoss = app_post_int('max_loss');
$chType = app_post_int('ch_type');

if ($moduleFirst === $moduleSecond) {
    $testType = 2;
}

app_stmt(
    'INSERT INTO `test_throughput` (`test_type`, `module_first`, `module_second`, `thr_begin`, `count`, `max_loss`, `ch_type`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, 1)',
    'iiiiiii',
    [$testType, $moduleFirst, $moduleSecond, $clock, $count, $maxLoss, $chType]
);

app_redirect('menu.php');
