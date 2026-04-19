<?php
require_once 'db.php';

app_require_post();

$moduleIp = app_post_string('moduleIP');
$pass = app_post_string('SNMPpassword');
$value = app_post_string('newvalue');
$oid = app_post_string('OID');
$type = app_post_string('type');

if ($moduleIp !== '' && $pass !== '' && $oid !== '' && $type !== '') {
    $session = new SNMP(SNMP::VERSION_2C, $moduleIp, $pass, 5000, 0);
    if ($session->get('.1.3.6.1.4.1.2010.1.1.0')) {
        $session->set($oid, $type, $value);
    }
    $session->close();
}

app_redirect('module_view_all.php');
