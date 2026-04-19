<?php
require_once 'db.php';

app_require_post();

$testType = app_post_int('test_type');
$testName = app_post_string('test_name');
$moduleFirst = app_post_int('module_first');
$moduleSecond = app_post_int('module_second');
$blockSize = app_post_int('block_size');
$clock = app_post_int('clock');
$count = app_post_int('count');
$delayMax = app_post_float('delay_max');
$jitterMax = app_post_float('jitter_max');
$delayOneMax = app_post_float('delay1_max');
$jitterOneMax = app_post_float('jitter1_max');
$lossMax = app_post_float('loss_max');

if ($moduleFirst === $moduleSecond) {
    $testType = 2;
}

$testDelay = app_post_bool('test_delay');
$testJitter = app_post_bool('test_jitter');
$testDelayOne = app_post_bool('test_delay_1');
$testJitterDelayOne = app_post_bool('test_jitter_delay_1');
$testLoss = app_post_bool('test_loss');

$countExisting = (int)app_scalar(
    'SELECT COUNT(*) FROM `test_sla_real` WHERE `module_first` = ? AND `module_second` = ?',
    'ii',
    [$moduleFirst, $moduleSecond]
);

if ($countExisting > 0) {
    app_redirect('add_test_real_double.php?' . http_build_query([
        'name' => $testName,
        'type' => $testType,
        'sfp1' => $moduleFirst,
        'sfp2' => $moduleSecond,
        'bsize' => $blockSize,
        'clock' => $clock,
        'count' => $count,
        'ch_delay' => $testDelay,
        'ch_jitter' => $testJitter,
        'ch_d1' => $testDelayOne,
        'ch_d1_j' => $testJitterDelayOne,
        'ch_loss' => $testLoss,
        'delay_max' => $delayMax,
        'jitter_max' => $jitterMax,
        'delay1_max' => $delayOneMax,
        'jitter1_max' => $jitterOneMax,
        'loss_max' => $lossMax,
    ]));
}

app_stmt(
    'INSERT INTO `test_sla_real` (`test_type`, `name`, `module_first`, `module_second`, `block_size`, `clock`, `count`, `status`, `test_delay`, `test_delay_1`, `test_delay1_jitter`, `test_delay_jitter`, `test_loss`, `delay_max`, `jitter_max`, `delay1_max`, `jitter1_max`, `loss_max`) VALUES (?, ?, ?, ?, ?, ?, ?, 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)',
    'isiiiiiiiiiiddddd',
    [$testType, $testName, $moduleFirst, $moduleSecond, $blockSize, $clock, $count, $testDelay, $testDelayOne, $testJitterDelayOne, $testJitter, $testLoss, $delayMax, $jitterMax, $delayOneMax, $jitterOneMax, $lossMax]
);

app_redirect('menu.php');
