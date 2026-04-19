<?php
require_once 'db.php';

app_require_post();

$testType = app_post_int('test_type');
$moduleFirst = app_post_int('module_first');
$moduleSecond = app_post_int('module_second');
$thrBegin = app_post_int('thr_begin');
$countProbPacks = app_post_int('count_prob_packs');
$countProbs = app_post_int('count_probs');

if ($moduleFirst === $moduleSecond) {
    $testType = 2;
}

app_stmt(
    'INSERT INTO `test_bert` (`test_type`, `module_first`, `module_second`, `thr_begin`, `count_prob_packs`, `count_probs`, `status`) VALUES (?, ?, ?, ?, ?, ?, 1)',
    'iiiiii',
    [$testType, $moduleFirst, $moduleSecond, $thrBegin, $countProbPacks, $countProbs]
);

app_redirect('menu.php');
