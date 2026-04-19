<?php
require_once 'db.php';

app_require_post();

$testType = app_post_int('test_type');
$moduleFirst = app_post_int('module_first');
$moduleSecond = app_post_int('module_second');
$clock = app_post_int('clock');
$countTests = app_post_int('count_tests');
$countPacks = app_post_int('count_packs');

if ($moduleFirst === $moduleSecond) {
    $testType = 2;
}

app_stmt(
    'INSERT INTO `test_latency` (`test_type`, `module_first`, `module_second`, `thr_begin`, `count_packs`, `count_tests`, `status`) VALUES (?, ?, ?, ?, ?, ?, 1)',
    'iiiiii',
    [$testType, $moduleFirst, $moduleSecond, $clock, $countPacks, $countTests]
);

app_redirect('menu.php');
