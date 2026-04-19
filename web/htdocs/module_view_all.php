<?php

require_once 'db.php';

if (isset($_SESSION['out_file'])) {
    unset($_SESSION['out_file']);
}

$link = app_db();
$sql = "select * from `modules_sfp_sla`";
$search = mysqli_query($link, $sql);
$number_sfp = mysqli_num_rows($search);

?>

<!DOCTYPE html>
<html class="no-js css-menubar" lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=0, minimal-ui">
    <meta name="description" content="SFP-SLA management console">
    <meta name="author" content="">
    <title>SFP-SLA Management Console</title>
    <link rel="apple-touch-icon" href="img/logo.svg">
    <link rel="shortcut icon" href="img/logo.svg">
    <!-- Stylesheets -->
    <link rel="stylesheet" href="css/bootstrap.min.css">
    <link rel="stylesheet" href="css/bootstrap-extend.min.css">
    <link rel="stylesheet" href="css/site.min.css">
    <!-- Plugins -->
    <link rel="stylesheet" href="vendor/animsition/animsition.css">
    <link rel="stylesheet" href="vendor/asscrollable/asScrollable.css">
    <link rel="stylesheet" href="vendor/waves/waves.css">
    <link rel="stylesheet" href="vendor/chartist-js/chartist.css">
    <link rel="stylesheet" href="vendor/chartist-plugin-tooltip/chartist-plugin-tooltip.css">
    <link rel="stylesheet" href="css/dashboard/v1.css">
    <!-- Fonts -->
    <link rel="stylesheet" href="fonts/material-design/material-design.min.css">
    <!-- Scripts -->
    <script src="vendor/modernizr/modernizr.js"></script>
    <script src="vendor/breakpoints/breakpoints.js"></script>
    <script>
        Breakpoints();
    </script>
</head>
<body class="site-navbar-small dashboard layout-boxed">

<nav class="site-navbar navbar navbar-inverse navbar-fixed-top navbar-expand-md navbar-mega navbar-inverse"
     role="navigation">
    <div class="navbar-header">

        <button type="button" class="navbar-toggle hamburger hamburger-close navbar-toggle-left hided"
                data-toggle="menubar">
            <span class="sr-only">Toggle navigation</span>
            <span class="hamburger-bar"></span>
        </button>

        <a class="navbar-brand navbar-brand-center" href="menu.php">
            <img class="navbar-brand-logo navbar-brand-logo-normal" src="img/logo.svg"
                 title="SLA-SFP">
            <span class="navbar-brand-text ">SFP-SLA</span>
        </a>


    </div>
    <div class="navbar-container container-fluid">
        <h3 class="navbar-text hidden-xs">SFP-SLA Management Console</h3>
    </div>
</nav>
<div class="site-menubar">
    <div class="site-menubar-body">
        <div>
            <div>
                <ul class="site-menu">
                    <li class="dropdown site-menu-item has-sub">
                        <a class="dropdown-toggle" href="javascript:void(0)" data-dropdown-toggle="false">
                            <i class="site-menu-icon md-view-compact" aria-hidden="true"></i>
                            <span class="site-menu-title">Global settings</span>
                            <span class="site-menu-arrow"></span>
                        </a>
                        <div class="dropdown-menu">
                            <div class="site-menu-scroll-wrap is-list">
                                <div>
                                    <div>
                                        <ul class="site-menu-sub site-menu-normal-list">
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="global_edite.php">
                                                    <span class="site-menu-title">Controller server</span>
                                                </a>
                                            </li>
</ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </li>
                    <li class="dropdown site-menu-item has-sub">
                        <a class="dropdown-toggle" href="javascript:void(0)" data-dropdown-toggle="false">
                            <i class="site-menu-icon md-google-pages" aria-hidden="true"></i>
                            <span class="site-menu-title">SFP-SLA modules</span>
                            <span class="site-menu-arrow"></span>
                        </a>
                        <div class="dropdown-menu">
                            <div class="site-menu-scroll-wrap is-list">
                                <div>
                                    <div>
                                        <ul class="site-menu-sub site-menu-normal-list">
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="module_view_all.php">
                                                    <span class="site-menu-title">View all</span>
                                                </a>
                                            </li>
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="module_add.php">
                                                    <span class="site-menu-title">Add module</span>
                                                </a>
                                            </li>
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="module_del.php">
                                                    <span class="site-menu-title">Delete module</span>
                                                </a>
                                            </li>

                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </li>

                    <li class="dropdown site-menu-item has-sub">
                        <a class="dropdown-toggle" href="javascript:void(0)" data-dropdown-toggle="false">
                            <i class="site-menu-icon md-apps" aria-hidden="true"></i>
                            <span class="site-menu-title">SLA tests</span>
                            <span class="site-menu-arrow"></span>
                        </a>
                        <div class="dropdown-menu">
                            <div class="site-menu-scroll-wrap is-list">
                                <div>
                                    <div>
                                        <ul class="site-menu-sub site-menu-normal-list">
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="tests_view_all.php">
                                                    <span class="site-menu-title">View all tests</span>
                                                </a>
                                            </li>
                                            <hr>
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="test_add_real.php">
                                                    <span class="site-menu-title">Add SLA test</span>
                                                </a>
                                            </li>
                                            <li class="site-menu-item ">
                                                <a class="animsition-link">
                                                    <span class="site-menu-title">RFC 2544 tests:</span>
                                                </a>
                                            </li>
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="test_add_throughput.php">
                                                    <span class="site-menu-title">-> Throughput test</span>
                                                </a>
                                            </li>
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="test_add_latency.php">
                                                    <span class="site-menu-title">-> Latency test</span>
                                                </a>
                                            </li>
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="test_add_frame_loss.php">
                                                    <span class="site-menu-title">-> Frame loss test</span>
                                                </a>
                                            </li>
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="test_add_bert.php">
                                                    <span class="site-menu-title">-> Burst test</span>
                                                </a>
                                            </li>
                                            <hr>
                                            <li class="site-menu-item">
                                                <a class="animsition-link" href="test_add_Y1564.php">
                                                    <span class="site-menu-title">Y.1564 test</span>
                                                </a>
                                            </li>
                                        </ul>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </li>
</ul>
            </div>
        </div>
    </div>
</div>
<!-- Page -->
<div class="page animsition">
    <div class="page-content container-fluid">
        <div class="row" data-plugin="matchHeight" data-by-row="true">


            <?php

            $sql = "select * from modules_sfp_sla";
            $search = mysqli_query($link, $sql);

            if ($search) {
                while ($row = mysqli_fetch_array($search)) {
                    $id = $row["id"];
                    $module_name = $row["name"];
                    $module_IP = $row["address_ip"];
                    $module_ver = $row["version"];
                    $module_loc = $row["location"];
                    $snmp_get = "";

                    $session_snmp = new SNMP(SNMP::VERSION_2C, $module_IP, "public", 5000);

                    if ($snmp_mac = $session_snmp->get(".1.3.6.1.4.1.2010.1.1.0")) {

                        $pos = strpos($snmp_mac, ':') + 1;
                        $snmp_mac = substr($snmp_mac, $pos, strlen($snmp_mac) - $pos);
                        $snmp_mac = dechex($snmp_mac);

                        $snmp_ip = $session_snmp->get(".1.3.6.1.4.1.2010.1.2.0");
                        $pos = strpos($snmp_ip, ':') + 1;
                        $snmp_ip = substr($snmp_ip, $pos, strlen($snmp_ip) - $pos);

                        $snmp_mask = $session_snmp->get(".1.3.6.1.4.1.2010.1.3.0");
                        $pos = strpos($snmp_mask, ':') + 1;
                        $snmp_mask = substr($snmp_mask, $pos, strlen($snmp_mask) - $pos);

                        $snmp_ip_gate = $session_snmp->get(".1.3.6.1.4.1.2010.1.4.0");
                        $pos = strpos($snmp_ip_gate, ':') + 1;
                        $snmp_ip_gate = substr($snmp_ip_gate, $pos, strlen($snmp_ip_gate) - $pos);

                        $snmp_ip_server_SLA = $session_snmp->get(".1.3.6.1.4.1.2010.1.5.0");
                        $pos = strpos($snmp_ip_server_SLA, ':') + 1;
                        $snmp_ip_server_SLA = substr($snmp_ip_server_SLA, $pos, strlen($snmp_ip_server_SLA) - $pos);

                        $snmp_ip_NTP = $session_snmp->get(".1.3.6.1.4.1.2010.1.6.0");
                        $pos = strpos($snmp_ip_NTP, ':') + 1;
                        $snmp_ip_NTP = substr($snmp_ip_NTP, $pos, strlen($snmp_ip_NTP) - $pos);

                        $snmp_adr_type = $session_snmp->get(".1.3.6.1.4.1.2010.1.7.0");
                        $snmp_write_eeprom = $session_snmp->get(".1.3.6.1.4.1.2010.1.8.0");
                        $snmp_time = $session_snmp->get(".1.3.6.1.4.1.2010.1.9.0");
                        $pos_b = strpos($snmp_time, '(') + 1;
                        $pos_e = strpos($snmp_time, ')');
                        $snmp_time = substr($snmp_time, $pos_b, $pos_e - $pos_b);
                        $snmp_time = getdate($snmp_time - (mktime(0, 0, 0, 1, 1, 1970) - mktime(0, 0, 0, 1, 1, 1900)));

                        $snmp_mac_NTP = $session_snmp->get(".1.3.6.1.4.1.2010.1.10.0");
                        $pos = strpos($snmp_mac_NTP, ':') + 1;
                        $snmp_mac_NTP = substr($snmp_mac_NTP, $pos, strlen($snmp_mac_NTP) - $pos);
                        $snmp_mac_NTP = dechex($snmp_mac_NTP);

                        $snmp_mac_gate = $session_snmp->get(".1.3.6.1.4.1.2010.1.11.0");
                        $pos = strpos($snmp_mac_gate, ':') + 1;
                        $snmp_mac_gate = substr($snmp_mac_gate, $pos, strlen($snmp_mac_gate) - $pos);
                        $snmp_mac_gate = dechex($snmp_mac_gate);

                        $snmp_mac_server_SLA = $session_snmp->get(".1.3.6.1.4.1.2010.1.12.0");
                        $pos = strpos($snmp_mac_server_SLA, ':') + 1;
                        $snmp_mac_server_SLA = substr($snmp_mac_server_SLA, $pos, strlen($snmp_mac_server_SLA) - $pos);
                        $snmp_mac_server_SLA = dechex($snmp_mac_server_SLA);

                        $snmp_load_to_laser = $session_snmp->get(".1.3.6.1.4.1.2010.1.13.0");
                        $pos = strpos($snmp_load_to_laser, ':') + 1;
                        $snmp_load_to_laser = substr($snmp_load_to_laser, $pos, strlen($snmp_load_to_laser) - $pos);

                        $snmp_load_to_comm = $session_snmp->get(".1.3.6.1.4.1.2010.1.14.0");
                        $pos = strpos($snmp_load_to_comm, ':') + 1;
                        $snmp_load_to_comm = substr($snmp_load_to_comm, $pos, strlen($snmp_load_to_comm) - $pos);

                        $snmp_NTP_period = $session_snmp->get(".1.3.6.1.4.1.2010.1.15.0");
                        $pos = strpos($snmp_NTP_period, ':') + 1;
                        $snmp_NTP_period = substr($snmp_NTP_period, $pos, strlen($snmp_NTP_period) - $pos);

                        $snmp_VLAN = $session_snmp->get(".1.3.6.1.4.1.2010.1.16.0");
                        $pos = strpos($snmp_VLAN, ':') + 1;
                        $snmp_VLAN = substr($snmp_VLAN, $pos, strlen($snmp_VLAN) - $pos);

                        $snmp_QinQ = $session_snmp->get(".1.3.6.1.4.1.2010.1.17.0");
                        $pos = strpos($snmp_QinQ, ':') + 1;
                        $snmp_QinQ = substr($snmp_QinQ, $pos, strlen($snmp_QinQ) - $pos);

                    } else {
                        $snmp_get = "SNMP request failed";
                    }
                    $session_snmp->close();


                    echo "
               <div class=\"row\">
                <div class=\"col-lg-12\">
                 <div class=\"accordion\" id=\"accordion$module_name\">
                    <div class=\"panel panel-bordered\">
                        <div class=\"panel-heading\">
                            <h3 class=\"panel-title\"><i class=\"icon md-link\" aria-hidden=\"true\"></i><button class=\"btn btn-link\" type=\"button\" data-toggle=\"collapse\" data-target=\"#collapseOne_$module_name\" aria-expanded=\"false\" aria-controls=\"collapseOne_$module_name\">Name: <strong>$module_name</strong> IP: <strong>$module_IP</strong></button></h3>
                        <div class=\"panel-actions\">
                                <form method=post action='action_del_module_go.php'>
                                    <button name=\"module_name\" value=\"$module_name\" type=\"submit\" class=\"btn btn-sm btn-icon btn-pure btn-default\"
                                            data-toggle=\"tooltip\" data-original-title=\"Delete\">
                                        <i class=\"icon md-close\" aria-hidden=\"true\"></i>
                                    </button>
                                </form>
                            </div>
                        
                        </div>
                        <div id=\"collapseOne_$module_name\" class=\"collapse\" aria-labelledby=\"headingOne\" data-parent=\"#accordion$module_name\">
                        <div class=\"panel-body\">
                        <div class=\"col-xlg-4 col-sm-12\">
                                
                                  <div> 
                                    <h5>  Firmware version: $module_ver </h5>
                                  </div>
                                  <div> 
                                    <h5>  Additional description: $module_loc </h5>
                                  </div>
                                  
                        </div>                         
                       <div class=\"col-xlg-8 col-sm-12\">
                                <div class=\"form-group form-material\">";
                    if ($snmp_get != "") {
                        echo "<h4 class=\"panel-footer bg-red-200\">Failed to read SNMP data</h4>
                            ";
                    } else {
                        echo "      <table class=\"table table-hover table-striped\">
                                        <thead>
                                        <tr>
                                            <th>SNMP parameter</th>
                                            <th>Value</th>
                                        </tr>
                                        </thead>
                                        <tbody>
                                        <tr>
                                        <td>SFP module MAC address</td>
                                        <td>$snmp_mac</td>
                                        </tr>
                                        <tr>
                                        <td>SFP module IP address</td>
                                        <td>$snmp_ip
                                        <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_IP_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                        
                                        
                                        </td>
                                        </tr>
                                        <tr>
                                        <td>SFP module IP mask</td>
                                        <td>$snmp_mask 
                                            <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_mask_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                        
                                        </td>
                                        </tr>
                                        <tr>
                                        <td>Gateway IP address</td>
                                        <td>$snmp_ip_gate
                                         <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_gate_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                        
                                        </td>
                                        </tr>
                                        <tr>
                                        <td>Gateway MAC address</td>
                                        <td>$snmp_mac_gate
                                         <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_gate_mac_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                         </td>
                                        </tr>
                                        
                                        <tr>
                                        <td>VLAN ID</td>
                                        <td>$snmp_VLAN
                                         <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_VLAN_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                        
                                        
</td>
                                        </tr>
                                        
                                        <tr>
                                        <td>VLAN ID QinQ</td>
                                        <td>$snmp_QinQ
                                         <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_QinQ_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                        
                                        </td>
                                        </tr>
                                        
                                        <tr>
                                            <td>SLA server IP address</td>
                                            <td>$snmp_ip_server_SLA
                                                <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_IP_server_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                            </td>
                                        </tr>
                                        <tr>
                                            <td>SLA server MAC address</td>
                                            <td>$snmp_mac_server_SLA
                                             
                                            </td>
                                        </tr>
                                        <tr>
                                        <td>Time server IP address</td>
                                            <td>$snmp_ip_NTP
                                                <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_IP_NTP_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                            </td>
                                        </tr>
                                        <tr>
                                        <td>Time server MAC address</td>
                                        <td>$snmp_mac_NTP</td>
                                        </tr>
                                        <td>Time synchronization period</td>
                                        <td>";
                        if ($snmp_NTP_period=='0'){
                            echo "disabled";
                        }
                        else {
                            echo pow(2, $snmp_NTP_period);
                            echo " sec.";
                        }
                        echo "
                                        <button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_NTP_period_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button>
                                        </td>
                                        </tr>
                                        <tr>
                                        <td>Addressing mode</td>";


                        if ($snmp_adr_type == '0') {
                            echo "<td>DHCP</td>";
                        } else {
                            echo "<td>Static addressing<button type=\"button\" class=\"btn btn-icon btn-flat pull-right\" data-toggle=\"modal\" data-target=\"#edit_type_SNMP_$snmp_mac\"><i class=\"icon md-edit\" aria-hidden=\"true\"></i></button></td>";
                        }
                        $snmp_year = $snmp_time['year'];
                        $snmp_mon = $snmp_time['mon'];
                        $snmp_mday = $snmp_time['mday'];
                        $snmp_hours = $snmp_time['hours'];
                        $snmp_minutes = $snmp_time['minutes'];
                        $snmp_seconds = $snmp_time['seconds'];
                        echo "</tr>
                                        <tr>
                                        <td>SFP system time</td>
                                        <td>$snmp_mday.$snmp_mon.$snmp_year   $snmp_hours:$snmp_minutes:$snmp_seconds</td>
                                        </tr>
                                        
                                        <tr>
                                        <td>Channel load toward laser</td>
                                         <td>";
                        $rez = ($snmp_load_to_laser * 8.0) / 1000.0;
                        echo "$rez Kbit/s</td>                                        
                                        </tr>
                                        <tr>
                                        <td>Channel load toward switch</td>
                                        <td>";
                        $rez = ($snmp_load_to_comm * 8.0) / 1000.0;
                        echo "$rez Kbit/s</td>
                                        </tr>
                                        </tbody>
                                        </table>";
                    }
                    echo "
                                </div>
                                </div>  
                                
                                <div class=\"text-right\">
                                     <form method=post action='modules_load_char_history.php'>   
                                       <button name=\"id\" value=\"$id\" type=\"submit\" class=\"btn btn-light waves-effect waves-light\">Show module channel load charts</button>                           
                                     </form>                            
                                </div>
                                
                                <div class=\"text-right\">
                                                             
                                       <button class='btn btn-success waves-effect waves-light' name='module_$snmp_mac' type='button' data-toggle=\"modal\" data-target=\"#edit_read_SNMP_$snmp_mac\" >Write to EEPROM</button>                           
                                                                  
                                </div>
                                
                                
                                </div>
                               </div>
                               </div>
                                
                                  ";

                    echo "
                       
                      
                <div class=\"modal fade\" id=\"edit_IP_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Change SFP module IP address</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
                            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.2.0\">
                            <input type=\"hidden\" name=\"type\" value=\"a\">
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">New value</label>
                                    <input class=\"form-control z-depth-1\" name=\"newvalue\"></input>
                                </div>
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\">Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                <!-- Modal -->       
                       
                        <!-- Modal -->
                <div class=\"modal fade\" id=\"edit_mask_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Change SFP module IP mask</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
                            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.3.0\">
                            <input type=\"hidden\" name=\"type\" value=\"a\">
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">New value</label>
                                    <input class=\"form-control z-depth-1\" name=\"newvalue\"></input>
                                </div>
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\"  >Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                <!-- Modal -->

                <!-- Modal -->
                <div class=\"modal fade\" id=\"edit_gate_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Change gateway IP address</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
                            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.4.0\">
                            <input type=\"hidden\" name=\"type\" value=\"a\">
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">New value</label>
                                    <input class=\"form-control z-depth-1\" name=\"newvalue\"></input>
                                </div>
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\">Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                <!-- Modal -->    
                
                  <!-- Modal -->
                <div class=\"modal fade\" id=\"edit_type_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Change addressing mode</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
                            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.7.0\">
                            <input type=\"hidden\" name=\"type\" value=\"i\">
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">New value</label>
                                    
                                    <select class=\"form-control z-depth-1\" name=\"newvalue\">
                                        <option value=\"0\">DHCP addressing</option>
                                        <option value=\"1\">Static addressing</option>
                                    </select>
                                    
                                </div>
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\">Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                <!-- Modal -->   

  <!-- Modal -->
                <div class=\"modal fade\" id=\"edit_VLAN_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Change VLAN number</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
                            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.16.0\">
                            <input type=\"hidden\" name=\"type\" value=\"i\">
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">New value (0 disables VLAN)</label>
                                    <input class=\"form-control z-depth-1\" name=\"newvalue\"></input>
                                </div>
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\">Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                <!-- Modal -->   

 <!-- Modal -->
                <div class=\"modal fade\" id=\"edit_IP_server_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Change SLA server IP address</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
                            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.5.0\">
                            <input type=\"hidden\" name=\"type\" value=\"a\">
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">New value</label>
                                    <input class=\"form-control z-depth-1\" name=\"newvalue\"></input>
                                </div>
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\">Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                <!-- Modal --> 
 <!-- Modal -->
                <div class=\"modal fade\" id=\"edit_IP_NTP_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Change time server IP address</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
                            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.6.0\">
                            <input type=\"hidden\" name=\"type\" value=\"a\">
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">New value</label>
                                    <input class=\"form-control z-depth-1\" name=\"newvalue\"></input>
                                </div>
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\">Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                <!-- Modal --> 
                
                <div class=\"modal fade\" id=\"edit_read_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Write value to EEPROM</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">       
                          
                           <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                           <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.8.0\">
                           <input type=\"hidden\" name=\"type\" value=\"i\">
                           <input type=\"hidden\" name=\"newvalue\" value=\"1\">
                           
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\">Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                
                
<!-- Modal -->
                <div class=\"modal fade\" id=\"edit_gate_mac_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"ModalLabel_$module_IP\" aria-hidden=\"true\">
                    <div class=\"modal-dialog\" role=\"document\">
                        <div class=\"modal-content\">
                            <div class=\"modal-header\">
                                <h5 class=\"modal-title\" id=\"ModalLabel_$module_IP\">Change gateway MAC address</h5>
                            </div>
                        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
                            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
                            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.11.0\">
                            <input type=\"hidden\" name=\"type\" value=\"Counter64\">
                            <div class=\"modal-body\">
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">New value</label>
                                    <input class=\"form-control z-depth-1\" name=\"newvalue\"></input>
                                </div>
                                <div class=\"form-group shadow-textarea\">
                                    <label for=\"FormControlTextarea\">Write password</label>
                                    <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
                                </div>
                            </div>
                            <div class=\"modal-footer\">
                                <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
                                <button type=\"submit\" class=\"btn btn-primary\">Save changes</button>
                            </div>
                        </form>
                        </div>
                    </div>
                </div>
                <!-- Modal -->   

<div class=\"modal fade\" id=\"edit_NTP_period_SNMP_$snmp_mac\" tabindex=\"-1\" role=\"dialog\" aria-labelledby=\"Modal_NTP_period_Label_$module_IP\" aria-hidden=\"true\">
  <div class=\"modal-dialog\" role=\"document\">
    <div class=\"modal-content\">
      <div class=\"modal-header\">
        <h5 class=\"modal-title\" id=\"Modal_NTP_period_Label_$module_IP\">Change NTP synchronization period</h5>

      </div>
        <form method=\"post\" action=\"action_edit_SNMP_go.php\">         
            <input type=\"hidden\" name=\"moduleIP\" value=\"$module_IP\">
            <input type=\"hidden\" name=\"OID\" value=\".1.3.6.1.4.1.2010.1.15.0\">
            <input type=\"hidden\" name=\"type\" value=\"i\">
      <div class=\"modal-body\">

          <div class=\"form-group shadow-textarea\">
              <label for=\"FormControlTextarea_NTP_period_\">New value (seconds)</label>
                          
              <select class=\"form-control z-depth-1\" name=\"newvalue\">
<option value=\"0\">do not update</option>
<option value=\"1\">2 seconds</option>
<option value=\"2\">4 seconds</option>
<option value=\"3\">8 seconds</option>
<option value=\"4\">16 seconds</option>
<option value=\"5\">32 seconds</option>
<option value=\"6\">64 seconds</option>
<option value=\"7\">2 minutes 8 seconds</option>
<option value=\"8\">4 minutes 16 seconds</option>
<option value=\"9\">8 minutes 32 seconds</option>
<option value=\"10\">17 minutes 4 seconds</option>
<option value=\"11\">34 minutes 8 seconds</option>
<option value=\"12\">1 hour 8 minutes 16 seconds</option>
<option value=\"13\">2 hours 16 minutes 32 seconds</option>
<option value=\"14\">4 hours 33 minutes 4 seconds</option>
<option value=\"15\">9 hours 6 minutes 8 seconds</option>
<option value=\"16\">18 hours 12 minutes 16 seconds</option>
<option value=\"17\">1 days 12 hours 24 minutes 32 seconds</option>
<option value=\"18\">3 days 49 minutes 4 seconds</option>
<option value=\"19\">6 days 1 hour 38 minutes 8 seconds</option>
<option value=\"20\">12 days 3 hours 16 minutes 16 seconds</option>
<option value=\"21\">2^21 seconds</option>
<option value=\"22\">2^22 seconds</option>
<option value=\"23\">2^23 seconds</option>
<option value=\"24\">2^24 seconds</option>
<option value=\"25\">2^25 seconds</option>
<option value=\"26\">2^26 seconds</option>
<option value=\"27\">2^27 seconds</option>
<option value=\"28\">2^28 seconds</option>
<option value=\"29\">2^29 seconds</option>
<option value=\"30\">2^30 seconds</option>
<option value=\"31\">2^31 seconds</option>
</select>
              
          </div>
          <div class=\"form-group shadow-textarea\">
              <label for=\"FormControlTextarea_NTP_period_\">Write password</label>
              <textarea class=\"form-control z-depth-1\" name=\"SNMPpassword\" rows=\"1\"></textarea>
          </div>

      </div>
      <div class=\"modal-footer\">
        <button type=\"reset\" class=\"btn btn-secondary\" data-dismiss=\"modal\">Cancel</button>
        <button type=\"submit\" class=\"btn btn-primary\" >Save changes</button>
    </div>
        </form>
  </div>
</div>
</div>
";

                }
                echo "
                    </div>
                    </div>
                    </div>";
            }
            
            ?>

        </div>
    </div>
</div>
<!-- End Page -->
<!-- Footer -->
<footer class="site-footer">
    <div class="site-footer-legal">© 2020 <a href="http://Server SFP SLA.ru/" target="_blank">Server SFP SLA</a></div>

</footer>
<!-- Core  -->
<script src="vendor/jquery/jquery.js"></script>
<script src="vendor/bootstrap/bootstrap.js"></script>
<script src="vendor/animsition/animsition.js"></script>
<script src="vendor/asscroll/jquery-asScroll.js"></script>
<script src="vendor/mousewheel/jquery.mousewheel.js"></script>
<script src="vendor/asscrollable/jquery.asScrollable.all.js"></script>
<script src="vendor/ashoverscroll/jquery-asHoverScroll.js"></script>
<script src="vendor/waves/waves.js"></script>
<!-- Plugins -->
<script src="vendor/chartist-js/chartist.min.js"></script>
<script src="vendor/chartist-plugin-tooltip/chartist-plugin-tooltip.min.js"></script>
<script src="vendor/matchheight/jquery.matchHeight-min.js"></script>
<!-- Scripts -->
<script src="js/core.js"></script>
<script src="js/site.js"></script>
<script src="js/sections/menu.js"></script>
<script src="js/sections/menubar.js"></script>
<script src="js/configs/config-colors.js"></script>
<script src="js/components/asscrollable.js"></script>
<script src="js/components/animsition.js"></script>
<script src="js/components/matchheight.js"></script>
<script src="js/dashboard/v1.js"></script>
</body>
</html>
