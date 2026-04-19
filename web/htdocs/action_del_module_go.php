<?php
require_once 'db.php';

app_require_post();

$moduleName = app_post_string('module_name');
if ($moduleName !== '') {
    app_stmt('DELETE FROM `modules_sfp_sla` WHERE `name` = ?', 's', [$moduleName]);
}

app_redirect('module_view_all.php');
