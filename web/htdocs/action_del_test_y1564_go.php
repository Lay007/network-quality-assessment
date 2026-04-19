<?php
require_once 'db.php';

app_require_post();

$id = app_post_int('id');
if ($id > 0) {
    app_stmt('DELETE FROM `test_y1564` WHERE `id` = ?', 'i', [$id]);
}

app_redirect('tests_view_all.php');
