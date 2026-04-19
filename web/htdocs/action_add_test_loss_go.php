<?php
require_once 'db.php';

app_require_post();

$testType = app_post_int('test_type');
$moduleFirst = app_post_int('module_first');
$moduleSecond = app_post_int('module_second');
$thrBegin = app_post_int('thr_begin');
$step = app_post_int('step');
$countFrames = app_post_int('count_frames');
$countSteps = app_post_int('count_steps');

if ($moduleFirst === $moduleSecond) {
    $testType = 2;
}

app_stmt(
    'INSERT INTO `test_frame_loss` (`test_type`, `module_first`, `module_second`, `thr_begin`, `step`, `count_frames`, `count_steps`, `status`) VALUES (?, ?, ?, ?, ?, ?, ?, 1)',
    'iiiiiii',
    [$testType, $moduleFirst, $moduleSecond, $thrBegin, $step, $countFrames, $countSteps]
);

app_redirect('menu.php');
