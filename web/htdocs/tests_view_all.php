<?php

require_once 'db.php';

if (isset($_SESSION['out_file'])) {
    unset($_SESSION['out_file']);
}

$link = app_db();

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


    <script src="/code/highstock.js"></script>
    <script src="/code/themes/grid-light.js"></script>
    <script src="/code/modules/exporting.js"></script>

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
            $link = app_db();

            $sql = "select * from test_sla_real order by id desc";
            // $sql = "select * from test_sla_real";
            $search = mysqli_query($link, $sql);

            if ($search) {
                while ($row = mysqli_fetch_array($search)) {
                    $id = (int)$row["id"];
                    $name = $row["name"];
                    $test_type = $row["test_type"];
                    $status = $row["status"];
                    $block_size = $row["block_size"];

                    $data_time_start = $row['data_start'];
                    $ch_delay = $row['test_delay'];
                    $ch_delay_jitter = $row['test_delay_jitter'];
                    $ch_delay1 = $row['test_delay_1'];
                    $ch_delay1_jitter = $row['test_delay1_jitter'];
                    $ch_loss = $row['test_loss'];

                    $delay_max = $row['delay_max'];
                    $delay1_max = $row['delay1_max'];

                    $jitter_max = $row['jitter_max'];
                    $jitter1_max = $row['jitter1_max'];
                    $loss_max = $row['loss_max'];

                    $per = $row['clock'];

                    $module_first = $row['module_first'];
                    $row_module = app_module_row((int)$module_first);
                    if ($row_module) {
                        $module1_name = $row_module['name'];
                        $module1_ip = $row_module['address_ip'];

                    }
                    $module_second = $row['module_second'];
                    $row_module = app_module_row((int)$module_second);
                    if ($row_module) {
                        $module2_name = $row_module['name'];
                        $module2_ip = $row_module['address_ip'];

                    }
                    $style_btn = "";
                    switch ($status) {
                        case 1:
                            $style_btn = "bg-yellow-100";
                            break;
                        case 2:
                            $style_btn = "bg-green-100";
                            break;
                        case 3:
                            $style_btn = "bg-blue-100";
                            break;
                        case 4:
                            $style_btn = "bg-red-100";
                            break;
                    }


                    echo "
              <div class=\"row\">
                <div class=\"col-lg-12\">
                <div class=\"accordion\" id=\"accordionTestReal$id\">
                   <div class=\"panel panel-bordered\">
                        <div class=\"panel-heading $style_btn\" >
                            <h3 class=\"panel-title\"><i class=\"icon md-chart\" aria-hidden=\"true\"></i><button class=\"btn btn-link\" type=\"button\" data-toggle=\"collapse\" data-target=\"#collapseTestReal_$id\" aria-expanded=\"false\" aria-controls=\"collapseTestReal_$id\"><h4>Real-time SLA parameter test</h4><h5>(SFP-SLA1: $module1_ip, SFP-SLA2: $module2_ip, started: $data_time_start)</h5></h3></button> 
                            <div class=\"panel-actions\">
                                <form method=post action='action_del_test_real_go.php'>
                                   <i class=\"icon\" aria-hidden=\"true\">
                                   <h5>";

                    switch ($status) {
                        case 1:
                            echo "Test is pending";
                            break;
                        case 2:
                            echo "Running";
                            break;
                        case 3:
                            echo "Test completed";
                            break;
                        case 4:
                            echo "Failed";
                            break;
                    }


                    echo "
                                   </h5>
                                   </i>
                                    <button name=\"id\" value=\"$id\" type=\"submit\" class=\"btn btn-sm btn-icon btn-pure btn-default\"
                                            data-toggle=\"tooltip\" data-original-title=\"Delete\">
                                        <i class=\"icon md-close\" aria-hidden=\"true\"></i>
                                    </button>
                                </form>
                            </div>
                        </div>
                      <div id=\"collapseTestReal_$id\" class=\"collapse\"  aria-labelledby=\"headingOne\" data-parent=\"#accordionTestReal$id\">
                      <div class=\"panel-body\">
                      
                       <div class=\"form-group\">
                     <table class=\"table table-hover table-striped\">
                            <thead>
                                        <tr>
                                            <th>Parameter</th>
                                            <th>Value</th>
                                        </tr>
                            </thead>
                            <tbody>
                            <tr>
                              <td>Test name</td>
                              <td>$name</td>
                            </tr>";

                    echo "<tr>                                        
                                                <td>Test scenario</td>
                                                <td>";
                    switch ($test_type) {
                        case 1:
                            echo "SFP-SLA1 - SFP-SLA2";
                            break;
                        case 2:
                            echo "Server - SFP-SLA1";
                            break;
                        case 3:
                            echo "Server - SFP-SLA2";
                            break;
                    }
                    echo "</td>
                                                </tr>

            
                                                <tr>
                    <td>First module SFP-SLA: </td>
                    <td>$module1_name : IP=  $module1_ip</td>
                            </tr>
                            
                            <tr>
                    <td>Second module SFP-SLA: </td>
                    <td>$module2_name : IP=  $module2_ip</td>
                            </tr>
                           <tr>
                             <td>Test frame size</td>
                              <td>$block_size</td>
                            </tr>
            
                           <tr>
                    <td>Packet generation period</td>
                    <td>$per</td>
                </tr>
                <tr>
                    <td>Measured parameters: </td>
                    <td>";

                    if ($ch_delay == 1) {
                        echo " Two-way delay";
                    }
                    if ($ch_delay_jitter == 1) {
                        echo ", two-way latency jitter";
                    }
                    if ($ch_delay1 == 1) {
                        echo ", one-way latency";
                    }
                    if ($ch_delay1_jitter == 1) {
                        echo ", one-way latency jitter";
                    }
                    if ($ch_loss == 1) {
                        echo ", packet loss ratio";
                    }

                    echo "
                    </td>
                </tr>";

                    if ($delay_max != 0) {
                        echo "
                            <tr>
                                <td>Maximum two-way delay (us):</td>
                                <td>$delay_max</td>
                            </tr>
                        ";
                    }
                    if ($jitter_max != 0) {
                        echo "
                            <tr>
                                <td>Maximum two-way jitter (us):</td>
                                <td>$jitter_max</td>
                            </tr>
                        ";
                    }
                    if ($delay1_max != 0) {
                        echo "
                            <tr>
                                <td>Maximum one-way delay (us):</td>
                                <td>$delay1_max</td>
                            </tr>
                        ";
                    }
                    if ($jitter1_max != 0) {
                        echo "
                            <tr>
                                <td>Maximum one-way jitter (us):</td>
                                <td>$jitter1_max</td>
                            </tr>
                        ";
                    }
                    if ($loss_max != 0) {
                        echo "
                            <tr>
                                <td>Maximum packet loss ratio:</td>
                                <td>$loss_max</td>
                            </tr>
                        ";
                    }
                    echo "
                                        </thead>
                                        <tbody>";


                    echo "<tr>                                        
                                                <td>Test scenario</td>
                                                <td>";
                    switch ($test_type) {
                        case 1:
                            echo "SFP-SLA1 - SFP-SLA2";
                            break;
                        case 2:
                            echo "Server - SFP-SLA1";
                            break;
                        case 3:
                            echo "Server - SFP-SLA2";
                            break;
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name1</td>
                                                <td>";

                    echo "$module1_name : IP=  $module1_ip";

                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name2</td>
                                                <td>";

                    echo "$module2_name : IP=  $module2_ip";

                    echo "</td>
                                                </tr>";

                    echo "<tr>                                        
                                                <td>Load change type</td>
                                                <td>";
                    switch ($ch_type) {
                        case 0:
                            echo "Binary search";
                            break;
                        case 1:
                            echo "Decrease by 10%";
                            break;
                        case 2:
                            echo "Decrease by 1%";
                            break;
                        case 3:
                            echo "Decrease by 0.01%";
                            break;
                    }
                    echo "</td></tr>";


                    echo "<tr>                                        
                                                <td>Max throughput [Mbit/s]</td>
                                                <td>$thr_begin</td>
                                                </tr>";

                    echo "<tr>                                        
                                                <td>Allowed error [%]</td>
                                                <td>$max_loss</td>
                                                </tr>";

                    echo "<tr>                                        
                                                <td>Test generation period (seconds)</td>
                                                <td>$count</td>
                                                </tr>
                                                <tr>
                                                <td>Test start time</td>
                                                <td>$data_time_start (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>                                                
                                                </tr>
                                                <tr>
                                                <td>Test completion time</td>
                                                <td>$data_time_end (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>
                                                </tr>
                                        </tbody>
                                    </table>";

                    switch ($status) {
                        case 1:
                            echo "<div class = \"panel-footer bg-yellow-100\">
                                                     <h4>Test is pending</h4>
                                                 </div>";
                            break;
                        case 2:
                            echo "<div class = \"panel-footer bg-green-100\">
                                                     <h4>Running</h4>
                                                 </div>";
                            break;
                        case 3:
                            echo "<div class = \"panel-footer bg-blue-100\">
                                                     <h4>Test completed</h4>
                                                 </div>";
                            break;
                        case 4:
                            echo "<div class = \"panel-footer bg-red-100\">
                                                     <h4>Test execution failed</h4>
                                                 </div>";
                            break;
                    }
                    echo "
            </div>
        </div>
        
      <div class=\"col-xlg-6 col-sm-12\">
        <!-- Panel Bar Stacked -->
        <div class=\"widget widget-shadow\" id=\"chartBarStacked_$id\">
            <div class=\"widget-content padding-30 height-full\">
                <div class=\"chart-header padding-bottom-10\" style=\"height:calc(100% - 350px);\">
                    <div class=\"font-size-20 inline-block\">Measurement result
                    </div>

                </div>
                <div class=\"ct-chart height-400\"></div>
            </div>
        </div>
    </div>";

                    /*    <script>
                            var Site = window.Site;
                            // widget chart
                            $(document).ready(function (jQuery) {
                                Site.run();
                                var stacked_bar = new Chartist.Bar('#chartBarStacked_$id .ct-chart', {
                                    labels: ['64', '128', '256', '512', '1024', '1280', '1518', '4096', '9000'],
                                    series: [";
                                        echo "[$rez_64, $rez_128, $rez_256, $rez_512, $rez_1024, $rez_1280, $rez_1518, $rez_4096, $rez_9000]";
                                        echo "
                                    ]
                                }, {
                                    stackBars: true,
                                    fullWidth: true,
                                    seriesBarDistance: 0,
                                    chartPadding: {
                                        top: 10,
                                        right: 0,
                                        bottom: 0,
                                        left: 0
                                    },
                                    axisX: {
                                        showLabel: true,
                                        showGrid: true,
                                        offset: 30
                                    },
                                    axisY: {
                                        showLabel: true,
                                        showGrid: true,
                                        offset: 80
                                    }
                                });

                            });

                        </script>
                    */
                    echo "
<script>
Highcharts.chart('chartBarStacked_$id', {
    chart: {
        type: 'column'
    },
    title: {
        text: 'Measurement results'
    },
    xAxis: {
        type: 'category',
        title: {
            text: 'Packet size (bytes)'
        },
        labels: {
            rotation: -45,
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },
    yAxis: {
        min: 0,
        title: {
            text: 'Throughput (Mbit/s)'
        }
    },
    legend: {
        enabled: false
    },
    tooltip: {
        pointFormat: 'Throughput: <b>{point.y:.1f} Mbit/s</b>'
    },
    credits: {
        enabled: false
    },
    series: [{
        name: 'Packet size',
        data: [
            ['64', $rez_64],
            ['128', $rez_128],
            ['256', $rez_256],
            ['512', $rez_512],
            ['1024', $rez_1024],
            ['1280', $rez_1280],
            ['1518', $rez_1518]
         //   ,
         //   ['4096', $rez_4096],
         //   ['9000', $rez_9000],

        ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.1f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    }]
});
</script>
</div>
 </div>
 </div>
</div>
</div>
</div>";

                }
            }

            $sql = "select * from test_latency order by id desc";
            $search = mysqli_query($link, $sql);
            if ($search) {
                while ($row = mysqli_fetch_array($search)) {
                    $id = (int)$row["id"];

                    $status = $row["status"];

                    $style_btn = "";
                    switch ($status) {
                        case 1:
                            $style_btn = "bg-yellow-100";
                            break;
                        case 2:
                            $style_btn = "bg-green-100";
                            break;
                        case 3:
                            $style_btn = "bg-blue-100";
                            break;
                        case 4:
                            $style_btn = "bg-red-100";
                            break;
                    }

                    $module_id_SFP1 = $row["module_first"];
                    $module_id_SFP2 = $row["module_second"];

                    $thr_begin = $row["thr_begin"];


                    $data_time_start = $row["datetime_start"];
                    $data_time_end = $row["datetime_end"];

                    echo "                 
            <div class=\"row\">
                <div class=\"col-lg-12\">    
                <div class=\"accordion\" id=\"accordionMain\">            
                    <div class=\"panel panel-bordered\">
                        <div class=\"panel-heading $style_btn\">                       
                            <div class=\"panel-title\"><i class=\"icon md-chart\" aria-hidden=\"true\"></i><button class=\"btn btn-flat waves-effect\" type=\"button\" data-toggle=\"collapse\" data-target=\"#collapse_delay_$id\" aria-expanded=\"false\" aria-controls=\"collapse_delay_$id\"><h4>Latency load test</h4></button></div> 

                            <div class=\"panel-actions\">
                                                       
                                <form method=post action='action_del_test_delay_go.php'>
                                    <button name=\"id\" value=\"$id\" type=\"submit\" class=\"btn btn-sm btn-icon btn-pure btn-default\"
                                            data-toggle=\"tooltip\" data-original-title=\"Delete\">
                                        <i class=\"icon md-close\" aria-hidden=\"true\"></i>
                                    </button>
                                </form>
                           
                            </div>
                        </div>
                        <div class=\"collapse\" id=\"collapse_delay_$id\">
                       

                        <div class=\"panel-body\" >
                            <div class=\"col-xlg-6 col-sm-12\">
                                <div class=\"form-group form-material\">

                                    <table class=\"table table-hover table-striped\">
                                        <thead>
                                        <tr>
                                            <th>Parameter</th>
                                            <th>Value</th>
                                        </tr>
                                        </thead>
                                        <tbody>";

                    $count = $row["count_packs"];
                    $test_type = $row["test_type"];


                    $rez_64 = $row["rez_64"];
                    $rez_64_max = $row["rez_64_max"] ;
                    $rez_64_min = $row["rez_64_min"] ;
                    $rez_128 = $row["rez_128"];
                    $rez_128_max = $row["rez_128_max"];
                    $rez_128_min = $row["rez_128_min"] ;
                    $rez_256 = $row["rez_256"];
                    $rez_256_max = $row["rez_256_max"];
                    $rez_256_min = $row["rez_256_min"];
                    $rez_512 = $row["rez_512"] ;
                    $rez_512_max = $row["rez_512_max"] ;
                    $rez_512_min = $row["rez_512_min"] ;
                    $rez_1024 = $row["rez_1024"] ;
                    $rez_1024_max = $row["rez_1024_max"];
                    $rez_1024_min = $row["rez_1024_min"];
                    $rez_1280 = $row["rez_1280"];
                    $rez_1280_max = $row["rez_1280_max"];
                    $rez_1280_min = $row["rez_1280_min"];
                    $rez_1518 = $row["rez_1518"] ;
                    $rez_1518_max = $row["rez_1518_max"] ;
                    $rez_1518_min = $row["rez_1518_min"];
                    //  $rez_4096 = $row["rez_4096"];
                    //  $rez_9000 = $row["rez_9000"];



                    echo "<tr>                                        
                                                <td>Test scenario</td>
                                                <td>";
                    switch ($test_type) {
                        case 1:
                            echo "SFP-SLA1 - SFP-SLA2";
                            break;
                        case 2:
                            echo "Server - SFP-SLA1";
                            break;
                        case 3:
                            echo "Server - SFP-SLA2";
                            break;
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name1</td>
                                                <td>";
                    $row_module = app_module_row((int)$module_id_SFP1);
                    if ($row_module) {
                        $module_name = $row_module['name'];
                        $module_ip = $row_module['address_ip'];
                        echo "$module_name : IP=  $module_ip";
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name2</td>
                                                <td>";
                    $row_module = app_module_row((int)$module_id_SFP2);
                    if ($row_module) {
                        $module_name = $row_module['name'];
                        $module_ip = $row_module['address_ip'];
                        echo "$module_name : IP=  $module_ip";
                    }
                    echo "</td>
                                                </tr>";


                    echo "<tr>                                        
                                                <td>Max throughput [Mbit/s]</td>
                                                <td>$thr_begin</td>
                                                </tr>";


                    echo "<tr>                                        
                                                <td>Test generation period (seconds)</td>
                                                <td>$count</td>
                                                </tr>
                                                <tr>
                                                <td>Test start time</td>
                                                <td>$data_time_start (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>
</tr>
                                                <tr>
                                                <td>Test completion time</td>
                                                <td>$data_time_end (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>
</tr>
                                        </tbody>
                                    </table>";

                    switch ($status) {
                        case 1:
                            echo "<div class = \"panel-footer bg-yellow-100\">
                                                     <h4>Test is pending</h4>
                                                 </div>";
                            break;
                        case 2:
                            echo "<div class = \"panel-footer bg-green-100\">
                                                     <h4>Running</h4>
                                                 </div>";
                            break;
                        case 3:
                            echo "<div class = \"panel-footer bg-blue-100\">
                                                     <h4>Test completed</h4>
                                                 </div>";
                            break;
                        case 4:
                            echo "<div class = \"panel-footer bg-red-100\">
                                                     <h4>Test execution failed</h4>
                                                 </div>";
                            break;
                    }
                    echo "
            </div>
        </div>
        
      <div class=\"col-xlg-6 col-sm-12\">
        <!-- Panel Bar Stacked -->
        <div class=\"widget widget-shadow\" id=\"chartBarDelay_$id\">
            <div class=\"widget-content padding-30 height-full\">
                <div class=\"chart-header padding-bottom-10\" style=\"height:calc(100% - 350px);\">
                    <div class=\"font-size-20 inline-block\">Measurement result
                    </div>

                </div>
                <div class=\"ct-chart height-400\"></div>
            </div>
        </div>
    </div>";


                    echo "
<script>
Highcharts.chart('chartBarDelay_$id', {    
    title: {
        text: 'Measurement results'
    },
    xAxis: {
        type: 'category',
        title: {
            text: 'Packet size (bytes)'
        },
        labels: {
            rotation: -45,
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },
    yAxis: {
        min: 0,
        title: {
            text: 'Latency value (us)'
        }
    },
     legend: {
        layout: 'vertical',
        align: 'right',
        verticalAlign: 'middle'
    },
    tooltip: {
        pointFormat: 'Latency value: <b>{point.y:.1f} us</b>'
    },
    credits: {
        enabled: false
    },
    series: [{
        name: 'Maximum latency',
        data: [
            ['64', $rez_64_max],
            ['128', $rez_128_max],
            ['256', $rez_256_max],
            ['512', $rez_512_max],
            ['1024', $rez_1024_max],
            ['1280', $rez_1280_max],
            ['1518', $rez_1518_max]     

        ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.1f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },
    {
        name: 'Average latency',
        data: [
            ['64', $rez_64],
            ['128', $rez_128],
            ['256', $rez_256],
            ['512', $rez_512],
            ['1024', $rez_1024],
            ['1280', $rez_1280],
            ['1518', $rez_1518]
       
        ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.1f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },
    {
        name: 'Minimum latency',
        data: [
            ['64', $rez_64_min],
            ['128', $rez_128_min],
            ['256', $rez_256_min],
            ['512', $rez_512_min],
            ['1024', $rez_1024_min],
            ['1280', $rez_1280_min],
            ['1518', $rez_1518_min]       

        ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.1f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    }]
});
</script>
</div>
 </div>
 </div>
</div>
</div>
</div>";

                }
            }

            $sql = "select * from test_frame_loss order by id desc";
            $search = mysqli_query($link, $sql);
            if ($search) {
                while ($row = mysqli_fetch_array($search)) {
                    $id = (int)$row["id"];

                    $count_steps = $row["count_steps"];
                    $step = $row["step"];
                    $count_frames = $row["count_frames"];

                    $test_type = $row["test_type"];
                    $status = $row["status"];

                    $module_id_SFP1 = $row["module_first"];
                    $module_id_SFP2 = $row["module_second"];

                    $thr_begin = $row["thr_begin"];


                    $data_time_start = $row["datetime_start"];
                    $data_time_end = $row["datetime_end"];

                    $style_btn = "";
                    switch ($status) {
                        case 1:
                            $style_btn = "bg-yellow-100";
                            break;
                        case 2:
                            $style_btn = "bg-green-100";
                            break;
                        case 3:
                            $style_btn = "bg-blue-100";
                            break;
                        case 4:
                            $style_btn = "bg-red-100";
                            break;
                    }

                    echo "                 
            <div class=\"row\">
                <div class=\"col-lg-12\">    
                <div class=\"accordion\" id=\"accordionMain\">            
                    <div class=\"panel panel-bordered\">
                        <div class=\"panel-heading $style_btn\">                       
                            <div class=\"panel-title\"><i class=\"icon md-chart\" aria-hidden=\"true\"></i><button class=\"btn btn-flat waves-effect\" type=\"button\" data-toggle=\"collapse\" data-target=\"#collapse_loss_$id\" aria-expanded=\"false\" aria-controls=\"collapse_loss_$id\"><h4>Frame loss load test</h4></button></div> 

                            <div class=\"panel-actions\">
                                                       
                                <form method=post action='action_del_test_loss_go.php'>
                                    <button name=\"id\" value=\"$id\" type=\"submit\" class=\"btn btn-sm btn-icon btn-pure btn-default\"
                                            data-toggle=\"tooltip\" data-original-title=\"Delete\">
                                        <i class=\"icon md-close\" aria-hidden=\"true\"></i>
                                    </button>
                                </form>
                           
                            </div>
                        </div>
                        <div class=\"collapse\" id=\"collapse_loss_$id\">
                       

                        <div class=\"panel-body\" >
                            <div class=\"col-xlg-6 col-sm-12\">
                                <div class=\"form-group form-material\">

                                    <table class=\"table table-hover table-striped\">
                                        <thead>
                                        <tr>
                                            <th>Parameter</th>
                                            <th>Value</th>
                                        </tr>
                                        </thead>
                                        <tbody>";

                    echo "<tr>                                        
                                                <td>Test scenario</td>
                                                <td>";
                    switch ($test_type) {
                        case 1:
                            echo "SFP-SLA1 - SFP-SLA2";
                            break;
                        case 2:
                            echo "Server - SFP-SLA1";
                            break;
                        case 3:
                            echo "Server - SFP-SLA2";
                            break;
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name1</td>
                                                <td>";
                    $row_module = app_module_row((int)$module_id_SFP1);
                    if ($row_module) {
                        $module_name = $row_module['name'];
                        $module_ip = $row_module['address_ip'];
                        echo "$module_name : IP=  $module_ip";
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name2</td>
                                                <td>";
                    $row_module = app_module_row((int)$module_id_SFP2);
                    if ($row_module) {
                        $module_name = $row_module['name'];
                        $module_ip = $row_module['address_ip'];
                        echo "$module_name : IP=  $module_ip";
                    }
                    echo "</td>
                                                </tr>";


                    echo "<tr>                                        
                                                <td>Max throughput [Mbit/s]</td>
                                                <td>$thr_begin</td>
                                                </tr>";


                    echo "<tr>                                        
                                                <td>Test generation period (seconds)</td>
                                                <td>$count_frames</td>
                                                </tr>
                                                
                                                                                 <tr>
                                                <td>Test start time</td>
                                                <td>$data_time_start (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>
</tr>
                                                
                                                <tr>
                                                <td>Test completion time</td>
                                                <td>$data_time_end (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>
</tr>
                                        </tbody>
                                    </table>";

                    switch ($status) {
                        case 1:
                            echo "<div class = \"panel-footer bg-yellow-100\">
                                                     <h4>Test is pending</h4>
                                                 </div>";
                            break;
                        case 2:
                            echo "<div class = \"panel-footer bg-green-100\">
                                                     <h4>Running</h4>
                                                 </div>";
                            break;
                        case 3:
                            echo "<div class = \"panel-footer bg-blue-100\">
                                                     <h4>Test completed</h4>
                                                 </div>";
                            break;
                        case 4:
                            echo "<div class = \"panel-footer bg-red-100\">
                                                     <h4>Test execution failed</h4>
                                                 </div>";
                            break;
                    }
                    echo "
            </div>
        </div>
        
      <div class=\"col-xlg-6 col-sm-12\">
        <!-- Panel Bar Stacked -->
        <div class=\"widget widget-shadow\" id=\"chartBarLoss_$id\">
            <div class=\"widget-content padding-30 height-full\">
                <div class=\"chart-header padding-bottom-10\" style=\"height:calc(100% - 350px);\">
                    <div class=\"font-size-20 inline-block\">Measurement result
                    </div>

                </div>
                <div class=\"ct-chart height-400\"></div>
            </div>
        </div>
    </div>";


                    echo "
<script>
Highcharts.chart('chartBarLoss_$id', {
    
    title: {
        text: 'Measurement results'
    },
    xAxis: {
        type: 'category',
        title: {
            text: 'Channel load (%)'
        },
        labels: {
            rotation: -45,
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },
    yAxis: {
        min: 0,
        title: {
            text: 'Packet loss ratio'
        }
    },
     legend: {
        layout: 'vertical',
        align: 'right',
        verticalAlign: 'middle'
    },
    tooltip: {
        pointFormat: 'Packet loss ratio: <b>{point.y:.5f}</b>'
    },
    credits: {
        enabled: false
    },
    series: [
        
        {           
        name: '64-byte packets',
        data: [";

                    $stmt_rez = app_stmt('select rez_64 from test_frame_loss_rez where id_test = ? order by id asc', 'i', [$id]);
                    $search_rez = mysqli_stmt_get_result($stmt_rez);
                    $row_rez = mysqli_fetch_array($search_rez);
                    $cc = 0;
                    do {
                        $rez = 100 - $step * $cc;
                        $cc++;
                        if ($cc > 1) {
                            printf(",['$rez'," . $row_rez ["rez_64"] . "]");
                        } else {
                            printf("['$rez'," . $row_rez ["rez_64"] . "]");
                        }

                    } while ($row_rez = mysqli_fetch_array($search_rez));

                    echo "    ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.4f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },";

                    echo "    {           
        name: '128-byte packets',
        data: [";

                    $stmt_rez = app_stmt('select rez_128 from test_frame_loss_rez where id_test = ? order by id asc', 'i', [$id]);
                    $search_rez = mysqli_stmt_get_result($stmt_rez);
                    $row_rez = mysqli_fetch_array($search_rez);
                    $cc = 0;
                    do {
                        $rez = 100 - $step * $cc;
                        $cc++;
                        if ($cc > 1) {
                            printf(",['$rez'," . $row_rez ["rez_128"] . "]");
                        } else {
                            printf("['$rez'," . $row_rez ["rez_128"] . "]");
                        }

                    } while ($row_rez = mysqli_fetch_array($search_rez));

                    echo "    ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.4f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },";

                    echo "{
                        name: '256-byte packets',
        data: [";

                    $stmt_rez = app_stmt('select rez_256 from test_frame_loss_rez where id_test = ? order by id asc', 'i', [$id]);
                    $search_rez = mysqli_stmt_get_result($stmt_rez);
                    $row_rez = mysqli_fetch_array($search_rez);
                    $cc = 0;
                    do {
                        $rez = 100 - $step * $cc;
                        $cc++;
                        if ($cc > 1) {
                            printf(",['$rez'," . $row_rez ["rez_256"] . "]");
                        } else {
                            printf("['$rez'," . $row_rez ["rez_256"] . "]");
                        }

                    } while ($row_rez = mysqli_fetch_array($search_rez));

                    echo "    ],
        dataLabels: {
                        enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.4f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                            fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },";

                    echo "    {           
        name: '512-byte packets',
        data: [";

                    $stmt_rez = app_stmt('select rez_512 from test_frame_loss_rez where id_test = ? order by id asc', 'i', [$id]);
                    $search_rez = mysqli_stmt_get_result($stmt_rez);
                    $row_rez = mysqli_fetch_array($search_rez);
                    $cc = 0;
                    do {
                        $rez = 100 - $step * $cc;
                        $cc++;
                        if ($cc > 1) {
                            printf(",['$rez'," . $row_rez ["rez_512"] . "]");
                        } else {
                            printf("['$rez'," . $row_rez ["rez_512"] . "]");
                        }

                    } while ($row_rez = mysqli_fetch_array($search_rez));

                    echo "    ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.4f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },";

                    echo "    {           
        name: '1024-byte packets',
        data: [";

                    $stmt_rez = app_stmt('select rez_1024 from test_frame_loss_rez where id_test = ? order by id asc', 'i', [$id]);
                    $search_rez = mysqli_stmt_get_result($stmt_rez);
                    $row_rez = mysqli_fetch_array($search_rez);
                    $cc = 0;
                    do {
                        $rez = 100 - $step * $cc;
                        $cc++;
                        if ($cc > 1) {
                            printf(",['$rez'," . $row_rez ["rez_1024"] . "]");
                        } else {
                            printf("['$rez'," . $row_rez ["rez_1024"] . "]");
                        }

                    } while ($row_rez = mysqli_fetch_array($search_rez));

                    echo "    ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.4f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },";

                    echo "    {           
        name: '1280-byte packets',
        data: [";

                    $stmt_rez = app_stmt('select rez_1280 from test_frame_loss_rez where id_test = ? order by id asc', 'i', [$id]);
                    $search_rez = mysqli_stmt_get_result($stmt_rez);
                    $row_rez = mysqli_fetch_array($search_rez);
                    $cc = 0;
                    do {
                        $rez = 100 - $step * $cc;
                        $cc++;
                        if ($cc > 1) {
                            printf(",['$rez'," . $row_rez ["rez_1280"] . "]");
                        } else {
                            printf("['$rez'," . $row_rez ["rez_1280"] . "]");
                        }

                    } while ($row_rez = mysqli_fetch_array($search_rez));

                    echo "    ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.4f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },";


                    echo "    {           
        name: '1518-byte packets',
        data: [";

                    $stmt_rez = app_stmt('select rez_1518 from test_frame_loss_rez where id_test = ? order by id asc', 'i', [$id]);
                    $search_rez = mysqli_stmt_get_result($stmt_rez);
                    $row_rez = mysqli_fetch_array($search_rez);
                    $cc = 0;
                    do {
                        $rez = 100 - $step * $cc;
                        $cc++;
                        if ($cc > 1) {
                            printf(",['$rez'," . $row_rez ["rez_1518"] . "]");
                        } else {
                            printf("['$rez'," . $row_rez ["rez_1518"] . "]");
                        }

                    } while ($row_rez = mysqli_fetch_array($search_rez));

                    echo "    ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.4f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    }]
});
</script>
</div>
 </div>
 </div>
</div>
</div>
</div>";

                }
            }

            $sql = "select * from test_bert order by id desc";
            $search = mysqli_query($link, $sql);
            if ($search) {
                while ($row = mysqli_fetch_array($search)) {
                    $id = (int)$row["id"];
                    $count_prop_packs = $row["count_prob_packs"];
                    $count_probs = $row["count_probs"];
                    $test_type = $row["test_type"];
                    $status = $row["status"];

                    $style_btn = "";
                    switch ($status) {
                        case 1:
                            $style_btn = "bg-yellow-100";
                            break;
                        case 2:
                            $style_btn = "bg-green-100";
                            break;
                        case 3:
                            $style_btn = "bg-blue-100";
                            break;
                        case 4:
                            $style_btn = "bg-red-100";
                            break;
                    }

                    $rez_64 = $row["rez_64"] / 1000.0;

                    $rez_128 = $row["rez_128"] / 1000.0;

                    $rez_256 = $row["rez_256"] / 1000.0;

                    $rez_512 = $row["rez_512"] / 1000.0;

                    $rez_1024 = $row["rez_1024"] / 1000.0;

                    $rez_1280 = $row["rez_1280"] / 1000.0;

                    $rez_1518 = $row["rez_1518"] / 1000.0;

                    //  $rez_4096 = $row["rez_4096"];
                    //  $rez_9000 = $row["rez_9000"];

                    $module_id_SFP1 = $row["module_first"];
                    $module_id_SFP2 = $row["module_second"];

                    $thr_begin = $row["thr_begin"];


                    $data_time_start = $row["datetime_start"];
                    $data_time_end = $row["datetime_end"];

                    echo "                 
            <div class=\"row\">
                <div class=\"col-lg-12\">    
                <div class=\"accordion\" id=\"accordionMain\">            
                    <div class=\"panel panel-bordered\">
                        <div class=\"panel-heading $style_btn\">                       
                            <div class=\"panel-title\"><i class=\"icon md-chart\" aria-hidden=\"true\"></i><button class=\"btn btn-flat waves-effect\" type=\"button\" data-toggle=\"collapse\" data-target=\"#collapse_berst_$id\" aria-expanded=\"false\" aria-controls=\"collapse_berst_$id\"><h4>Burst throughput load test</h4></button></div> 

                            <div class=\"panel-actions\">
                                                       
                                <form method=post action='action_del_test_bert_go.php'>
                                    <button name=\"id\" value=\"$id\" type=\"submit\" class=\"btn btn-sm btn-icon btn-pure btn-default\"
                                            data-toggle=\"tooltip\" data-original-title=\"Delete\">
                                        <i class=\"icon md-close\" aria-hidden=\"true\"></i>
                                    </button>
                                </form>
                           
                            </div>
                        </div>
                        <div class=\"collapse\" id=\"collapse_berst_$id\">
                       

                        <div class=\"panel-body\" >
                            <div class=\"col-xlg-6 col-sm-12\">
                                <div class=\"form-group form-material\">

                                    <table class=\"table table-hover table-striped\">
                                        <thead>
                                        <tr>
                                            <th>Parameter</th>
                                            <th>Value</th>
                                        </tr>
                                        </thead>
                                        <tbody>";


                    echo "<tr>                                        
                                                <td>Test scenario</td>
                                                <td>";
                    switch ($test_type) {
                        case 1:
                            echo "SFP-SLA1 - SFP-SLA2";
                            break;
                        case 2:
                            echo "Server - SFP-SLA1";
                            break;
                        case 3:
                            echo "Server - SFP-SLA2";
                            break;
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name1</td>
                                                <td>";
                    $row_module = app_module_row((int)$module_id_SFP1);
                    if ($row_module) {
                        $module_name = $row_module['name'];
                        $module_ip = $row_module['address_ip'];
                        echo "$module_name : IP=  $module_ip";
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name2</td>
                                                <td>";
                    $row_module = app_module_row((int)$module_id_SFP2);
                    if ($row_module) {
                        $module_name = $row_module['name'];
                        $module_ip = $row_module['address_ip'];
                        echo "$module_name : IP=  $module_ip";
                    }
                    echo "</td>
                                                </tr>";


                    echo "<tr>                                        
                                                <td>Max throughput [Mbit/s]</td>
                                                <td>$thr_begin</td>
                                                </tr>";


                    echo "<tr>                                        
                                                <td>Test generation period (seconds)</td>
                                                <td>$count</td>
                                                </tr>
                                                
                                                <tr>                                        
                                                  <td>Number of test generations</td>
                                                  <td>$count_probs</td>
                                                </tr>
                                                
                                                <tr>
                                                <td>Test start time</td>
                                                <td>$data_time_start (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>
</tr>
<tr>
                                                <td>Test completion time</td>
                                                <td>$data_time_end (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>
</tr>
                                        </tbody>
                                    </table>";

                    switch ($status) {
                        case 1:
                            echo "<div class = \"panel-footer bg-yellow-100\">
                                                     <h4>Test is pending</h4>
                                                 </div>";
                            break;
                        case 2:
                            echo "<div class = \"panel-footer bg-green-100\">
                                                     <h4>Running</h4>
                                                 </div>";
                            break;
                        case 3:
                            echo "<div class = \"panel-footer bg-blue-100\">
                                                     <h4>Test completed</h4>
                                                 </div>";
                            break;
                        case 4:
                            echo "<div class = \"panel-footer bg-red-100\">
                                                     <h4>Test execution failed</h4>
                                                 </div>";
                            break;
                    }
                    echo "
            </div>
        </div>
        
      <div class=\"col-xlg-6 col-sm-12\">
        <!-- Panel Bar Stacked -->
        <div class=\"widget widget-shadow\" id=\"chartBarBerst_$id\">
            <div class=\"widget-content padding-30 height-full\">
                <div class=\"chart-header padding-bottom-10\" style=\"height:calc(100% - 350px);\">
                    <div class=\"font-size-20 inline-block\">Measurement result
                    </div>

                </div>
                <div class=\"ct-chart height-400\"></div>
            </div>
        </div>
    </div>";

                    echo "
<script>
Highcharts.chart('chartBarBerst_$id', {
    chart: {
        type: 'column'
    },
    title: {
        text: 'Measurement results'
    },
    xAxis: {
        type: 'category',
        title: {
            text: 'Packet size (bytes)'
        },
        labels: {
            rotation: -45,
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },
    yAxis: {
        min: 0,
        title: {
            text: 'Throughput (Mbit/s)'
        }
    },
    legend: {
        enabled: false
    },
    tooltip: {
        pointFormat: 'Throughput: <b>{point.y:.1f} Mbit/s</b>'
    },
    credits: {
        enabled: false
    },
    series: [{
        name: 'Packet size',
        data: [
            ['64', $rez_64],
            ['128', $rez_128],
            ['256', $rez_256],
            ['512', $rez_512],
            ['1024', $rez_1024],
            ['1280', $rez_1280],
            ['1518', $rez_1518]
         

        ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.1f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    }]
});
</script>
</div>
 </div>
 </div>
</div>
</div>
</div>";


                }
            }

            $sql = "select * from test_y1564 order by id desc";
            $search = mysqli_query($link, $sql);
            if ($search) {
                while ($row = mysqli_fetch_array($search)) {
                    $id = (int)$row["id"];

                    $test_type = $row["test_type"];
                    $status = $row["status"];

                    $style_btn = "";
                    switch ($status) {
                        case 1:
                            $style_btn = "bg-yellow-100";
                            break;
                        case 2:
                            $style_btn = "bg-green-100";
                            break;
                        case 3:
                            $style_btn = "bg-blue-100";
                            break;
                        case 4:
                            $style_btn = "bg-red-100";
                            break;
                    }

                    $rez_IR_s1 = $row["rez_IR_s1"];
                    $rez_FTD_s1 = $row["rez_FTD_s1"];
                    $rez_FVD_s1 = $row["rez_FVD_s1"];
                    $rez_FLR_s1 = $row["rez_FLR_s1"];

                    $rez_IR_s2 = $row["rez_IR_s2"];
                    $rez_FTD_s2 = $row["rez_FTD_s2"];
                    $rez_FVD_s2 = $row["rez_FVD_s2"];
                    $rez_FLR_s2 = $row["rez_FLR_s2"];

                    $rez_IR_s3 = $row["rez_IR_s3"];
                    $rez_FTD_s3 = $row["rez_FTD_s3"];
                    $rez_FVD_s3 = $row["rez_FVD_s3"];
                    $rez_FLR_s3 = $row["rez_FLR_s3"];

                    $rez_IR_s4 = $row["rez_IR_s4"];
                    $rez_FTD_s4 = $row["rez_FTD_s4"];
                    $rez_FVD_s4 = $row["rez_FVD_s4"];
                    $rez_FLR_s4 = $row["rez_FLR_s4"];

                    $rez_IR_eir = $row["rez_IR_eir"];
                    $rez_FTD_eir = $row["rez_FTD_eir"];
                    $rez_FVD_eir = $row["rez_FVD_eir"];
                    $rez_FLR_eir = $row["rez_FLR_eir"];

                    $rez_IR_tp = $row["rez_IR_tp"];
                    $rez_FTD_tp = $row["rez_FTD_tp"];
                    $rez_FVD_tp = $row["rez_FVD_tp"];
                    $rez_FLR_tp = $row["rez_FLR_tp"];

                    $module_id_SFP1 = $row["module_first"];
                    $module_id_SFP2 = $row["module_second"];

                    $block_size=$row["block_size"];
                    $ToS = $row["ToS"];
                    $VLAN_priority = $row["VLAN_priority"];

                    $step_count = $row["step_count"];
                    $period = $row["period"];

                    $CIR = $row["CIR"];
                    $EIR = $row["EIR"];
                    $TP = $row["TP"];

                    $graphEID = $CIR+$EIR;
                    $graphTP = $CIR+$TP;

                    $graph1st = $CIR*0.5;
                    $graph2st = $CIR*0.75;
                    $graph3st = $CIR*0.9;

                    $data_time_start = $row["datetime_start"];
                    $data_time_end = $row["datetime_end"];
                    echo "                 
            <div class=\"row\">
                <div class=\"col-lg-12\">    
                <div class=\"accordion\" id=\"accordionMain\">            
                    <div class=\"panel panel-bordered\">
                        <div class=\"panel-heading $style_btn\">                       
                            <div class=\"panel-title\"><i class=\"icon md-chart\" aria-hidden=\"true\"></i><button class=\"btn btn-flat waves-effect\" type=\"button\" data-toggle=\"collapse\" data-target=\"#collapse_Y1564_$id\" aria-expanded=\"false\" aria-controls=\"collapse_Y1564_$id\"><h4>Y.1564 load test</h4></button></div> 

                            <div class=\"panel-actions\">
                                                       
                                <form method=post action='action_del_test_y1564_go.php'>
                                    <button name=\"id\" value=\"$id\" type=\"submit\" class=\"btn btn-sm btn-icon btn-pure btn-default\"
                                            data-toggle=\"tooltip\" data-original-title=\"Delete\">
                                        <i class=\"icon md-close\" aria-hidden=\"true\"></i>
                                    </button>
                                </form>
                           
                            </div>
                        </div>
                        <div class=\"collapse\" id=\"collapse_Y1564_$id\">
                       

                        <div class=\"panel-body\" >
                            <div class=\"col-xlg-6 col-sm-12\">
                                <div class=\"form-group form-material\">

                                    <table class=\"table table-hover table-striped\">
                                        <thead>
                                        <tr>
                                            <th>Parameter</th>
                                            <th>Value</th>
                                        </tr>
                                        </thead>
                                        <tbody>";

                    echo "<tr>                                        
                                                <td>Test scenario</td>
                                                <td>";
                    switch ($test_type) {
                        case 1:
                            echo "SFP-SLA1 - SFP-SLA2";
                            break;
                        case 2:
                            echo "Server - SFP-SLA1";
                            break;
                        case 3:
                            echo "Server - SFP-SLA2";
                            break;
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name1</td>
                                                <td>";
                    $row_module = app_module_row((int)$module_id_SFP1);
                    if ($row_module) {
                        $module_name = $row_module['name'];
                        $module_ip = $row_module['address_ip'];
                        echo "$module_name : IP=  $module_ip";
                    }
                    echo "</td>
                                                </tr>";

                    echo "<tr>
                                                <td>SFP-SLA module name2</td>
                                                <td>";
                    $row_module = app_module_row((int)$module_id_SFP2);
                    if ($row_module) {
                        $module_name = $row_module['name'];
                        $module_ip = $row_module['address_ip'];
                        echo "$module_name : IP=  $module_ip";
                    }
                    echo "</td>
                                                </tr>";

echo " <tr><td>IP packet ToS field value</td>
            <td>";
                    switch ($ToS){
                         case '0':
                             echo "0000 - Routine service";
                             break;
                         case '1' :
                             echo "1000 - Minimum delay";
                             break;
                         case '2':
                             echo "0100 - Maximum throughput";
                             break;
                         case '4':
                             echo "0010 - Maximum reliability";
                             break;
                         case '8':
                             echo "0001 - Minimum cost";

                    }
                echo " </td>                         
                </tr>";

                    echo " <tr><td>VLAN Priority field value</td>
            <td>";
                    switch ($VLAN_priority){
                        case '0':
                            echo "000 - Routine service (Best Efforts)";
                            break;
                        case '1' :
                            echo "001 - Background traffic (Background)";
                            break;
                        case '2':
                            echo "010 - Excellent effort (Excellent Effort)";
                            break;
                        case '3':
                            echo "011 - Critical applications (Critical Applications)";
                            break;
                        case '4':
                            echo "100 - Video traffic";
                            break;
                        case '5':
                            echo "101 - Voice traffic";
                            break;
                        case '6' :
                            echo "110 - Internetwork control (Internetwork Control)";
                            break;
                        case '7':
                            echo "111 - Network control (Network Control)";
                            break;
                     }
                    echo " </td>                         
                </tr>
                
                <tr><td>Test packet size [bytes]</td><td>$block_size</td></tr>
                ";



                    echo "<tr>                                        
                                                <td>Max throughput [Mbit/s]</td>
                                                <td>$CIR</td>
                                                </tr>";


                    echo "                      <tr>                                        
                                                <td>Test generation period (seconds)</td>
                                                <td>$period</td>
                                                </tr>
                                                <tr>                                        
                                                <td>Step count</td>
                                                <td>$step_count</td>
                                                </tr>
                                                <tr>
                                                <td>Test start time</td>
                                                <td>$data_time_start (";
                                                    echo ini_get('date.timezone');
                                                    echo ") </td>
</tr>
                                                <tr>
                                                <td>Test completion time</td>
                                                <td>$data_time_end (";
                                                    echo ini_get('date.timezone');
                                                    echo ")</td>
</tr>
                                        </tbody>
                                    </table>";

                    switch ($status) {
                        case 1:
                            echo "<div class = \"panel-footer bg-yellow-100\">
                                                     <h4>Test is pending</h4>
                                                 </div>";
                            break;
                        case 2:
                            echo "<div class = \"panel-footer bg-green-100\">
                                                     <h4>Running</h4>
                                                 </div>";
                            break;
                        case 3:
                            echo "<div class = \"panel-footer bg-blue-100\">
                                                     <h4>Test completed</h4>
                                                 </div>";
                            break;
                        case 4:
                            echo "<div class = \"panel-footer bg-red-100\">
                                                     <h4>Test execution failed</h4>
                                                 </div>";
                            break;
                    }
                    echo "
            </div>
        </div>
        
      <div class=\"col-xlg-6 col-sm-12\">
        <!-- Panel Bar Stacked -->
        <div class=\"widget widget-shadow\" id=\"chartBarY1564_$id\">
            <div class=\"widget-content padding-30 height-full\">
                <div class=\"chart-header padding-bottom-10\" style=\"height:calc(100% - 350px);\">
                    <div class=\"font-size-20 inline-block\">Measurement result
                    </div>

                </div>
                <div class=\"ct-chart height-500\"></div>
            </div>
        </div>
    </div>
    
    
   <div class=\"col-xlg-12 col-sm-12\">
      <div class=\"form-group form-material\">
         <table class=\"table table-bordered table-hover table-striped\">
             <thead>
               <tr>
                 <th class=\"text-center\">Test step</th>
                 <th class=\"text-center\">(IR)<br>Throughput [Mbit/s]</th> 
                 <th class=\"text-center\">(FLR)<br>Packet loss ratio</th>
                 <th class=\"text-center\">(FTD)<br>Two-way latency [us]</th>
                 <th class=\"text-center\">(FDV)<br>Two-way latency jitter [us]</th>
               </tr>
             </thead>
             <tbody>
                <tr> <td colspan=\"5\" class=\"text-center\">Committed information rate test (CIR) </td></tr>               
                <tr>                                        
                   <td>Step 1 (50% CIR)</td>
                   <td class=\"text-center\">$rez_IR_s1</td>
                   <td class=\"text-center\">$rez_FLR_s1</td>
                   <td class=\"text-center\">$rez_FTD_s1</td>
                   <td class=\"text-center\">$rez_FVD_s1</td>
                </tr>   
                <tr>                
                   <td>Step 2 (75% CIR)</td>
                   <td class=\"text-center\">$rez_IR_s2</td>
                   <td class=\"text-center\">$rez_FLR_s2</td>
                   <td class=\"text-center\">$rez_FTD_s2</td>
                   <td class=\"text-center\">$rez_FVD_s2</td>
                </tr>
                <tr>   
                   <td>Step 3 (90% CIR)</td>
                   <td class=\"text-center\">$rez_IR_s3</td>
                   <td class=\"text-center\">$rez_FLR_s3</td>
                   <td class=\"text-center\">$rez_FTD_s3</td>
                   <td class=\"text-center\">$rez_FVD_s3</td>
                </tr>
                <tr>   
                   <td>Step 4 (100% CIR)</td> 
                   <td class=\"text-center\">$rez_IR_s4</td>
                   <td class=\"text-center\">$rez_FLR_s4</td>
                   <td class=\"text-center\">$rez_FTD_s4</td>
                   <td class=\"text-center\">$rez_FVD_s4</td>                  
                </tr>
                <tr> <td colspan=\"5\" class=\"text-center\">Excess information rate test (EIR) </td></tr> 
                <tr>   
                   <td>EIR test</td> 
                   <td class=\"text-center\">$rez_IR_eir</td>
                   <td class=\"text-center\">$rez_FLR_eir</td>
                   <td class=\"text-center\">$rez_FTD_eir</td>
                   <td class=\"text-center\">$rez_FVD_eir</td>                  
                </tr>
                <tr> <td colspan=\"5\" class=\"text-center\">Traffic policing test (TP) </td></tr> 
                <tr>   
                   <td>TP test</td> 
                   <td class=\"text-center\">$rez_IR_tp</td>
                   <td class=\"text-center\">$rez_FLR_tp</td>
                   <td class=\"text-center\">$rez_FTD_tp</td>
                   <td class=\"text-center\">$rez_FVD_tp</td>                  
                </tr>
             </tbody>
         </table>
      </div>
   </div>";




                    /*    <script>
                            var Site = window.Site;
                            // widget chart
                            $(document).ready(function (jQuery) {
                                Site.run();
                                var stacked_bar = new Chartist.Bar('#chartBarStacked_$id .ct-chart', {
                                    labels: ['64', '128', '256', '512', '1024', '1280', '1518', '4096', '9000'],
                                    series: [";
                                        echo "[$rez_64, $rez_128, $rez_256, $rez_512, $rez_1024, $rez_1280, $rez_1518, $rez_4096, $rez_9000]";
                                        echo "
                                    ]
                                }, {
                                    stackBars: true,
                                    fullWidth: true,
                                    seriesBarDistance: 0,
                                    chartPadding: {
                                        top: 10,
                                        right: 0,
                                        bottom: 0,
                                        left: 0
                                    },
                                    axisX: {
                                        showLabel: true,
                                        showGrid: true,
                                        offset: 30
                                    },
                                    axisY: {
                                        showLabel: true,
                                        showGrid: true,
                                        offset: 80
                                    }
                                });

                            });

                        </script>
                    */
                    echo "
<script>
Highcharts.chart('chartBarY1564_$id', {
    chart: {
        type: 'column'
    },
    title: {
        text: 'Measurement progress'
    },
    xAxis: {
        type: 'category',
        title: {
            text: 'Test step'
        },
        labels: {
            rotation: -45,
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }
    },
    yAxis: {
        min: 0,
        max: $graphTP,
        title: {
            text: 'Throughput (Mbit/s)'
        },
         plotBands: [{ 
            from: 0,
            to: $CIR,
            color: 'rgba(21, 192, 21, 0.2)',
            label: {
                text: 'CIR - committed information rate',
                style: {
                    color: '#606060'
                }
            }
        }, { // Light breeze
            from: $CIR,
            to: $graphEID,
            color: 'rgba(255, 215, 0, 0.3)',
            label: {
                text: 'Maximum CIR excess',
                style: {
                    color: '#606060'
                }
            }
        }, { // Gentle breeze
            from: $graphEID,
            to: $graphTP,
            color: 'rgba(255, 0,55, 0.3)',
            label: {
                text: 'Maximum TP excess',
                style: {
                    color: '#606060'
                }
            }
        }],
    },
    legend: {
        enabled: false
    },
    tooltip: {
        pointFormat: 'Throughput: <b>{point.y:.1f} Mbit/s</b>'
    },
    credits: {
        enabled: false
    },
    
    series: [{
        pointWidth: 80,
        color: 'rgba(0, 0,255, 0.4)',
        name: 'Throughput',
        data: [
            ['Step 1', $graph1st],
            ['Step 2', $graph2st],
            ['Step 3', $graph3st],
            ['Step 4', $CIR],
            ['EIR', $graphEID],
            ['TP', $graphTP]            
       ],
        dataLabels: {
            enabled: true,
            rotation: -90,
            color: '#FFFFFF',
            align: 'right',
            format: '{point.y:.1f}', // one decimal
            y: 10, // 10 pixels down from the top
            style: {
                fontSize: '13px',
                fontFamily: 'Verdana, sans-serif'
            }
        }        
    }]
});
</script>
</div>
 </div>
 </div>
</div>
</div>
</div>";

                }
            }


            
            ?>

        </div>
    </div>
</div>
<!-- End Page -->
<!-- Footer -->
<footer class="site-footer">
    <div class="site-footer-legal">© 2019 <a href="http://Server SFP SLA.ru/" target="_blank">Server SFP SLA</a></div>

</footer>
<!-- Core  -->
<script>
    $(document).ready(function () {
        $('.mdb-select').materialSelect();
    });
</script>

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
