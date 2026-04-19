<?php
require_once 'db.php';

app_require_post();

$serverIp = app_post_string('server_ip');
$netInterfaceName = app_post_string('net_interface_name');
$vlan = app_post_bool('vlan_check');
$vlanNumber = app_post_int('vlan_number', 0);
$qinq = app_post_bool('QinQ_check');
$qinqNumber = app_post_int('QinQ_number', 0);

app_stmt('TRUNCATE TABLE `global_config`');
app_stmt(
    'INSERT INTO `global_config` (`server_IP`, `net_interface_name`, `VLAN`, `VLAN_number`, `QinQ`, `QinQ_number`) VALUES (?, ?, ?, ?, ?, ?)',
    'ssiiii',
    [$serverIp, $netInterfaceName, $vlan, $vlanNumber, $qinq, $qinqNumber]
);

app_redirect('menu.php');
