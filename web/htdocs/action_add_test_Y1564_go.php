<?php
require_once 'db.php';

app_require_post();

$testType = app_post_int('test_type');
$moduleFirst = app_post_int('module_first');
$moduleSecond = app_post_int('module_second');
$blockSize = app_post_int('block_size');
$tos = app_post_int('ToS');
$vlanPriority = app_post_int('VLAN_priority', 0);
$cir = app_post_int('CIR');
$eir = app_post_int('EIR');
$tp = app_post_int('TP');
$period = app_post_int('period');
$stepCount = app_post_int('step_count');
$maxFtd = app_post_float('max_FTD');
$maxFvd = app_post_float('max_FVD');
$maxFlr = app_post_float('max_FLR');

if ($moduleFirst === $moduleSecond) {
    $testType = 2;
}

app_stmt(
    'INSERT INTO `test_y1564` (`test_type`, `module_first`, `module_second`, `block_size`, `ToS`, `VLAN_priority`, `CIR`, `EIR`, `TP`, `max_FTD`, `max_FVD`, `max_FLR`, `period`, `step_count`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1)',
    'iiiiiiiiidddii',
    [$testType, $moduleFirst, $moduleSecond, $blockSize, $tos, $vlanPriority, $cir, $eir, $tp, $maxFtd, $maxFvd, $maxFlr, $period, $stepCount]
);

app_redirect('menu.php');
