<?php

require_once 'db.php';

if (isset($_SESSION['out_file'])) {
    unset($_SESSION['out_file']);
}

$link = app_db();
$moduleId = app_post_int('id');



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
    <!--!
        <script src="/code/highcharts.js"></script>
        <script src="/code/highcharts-more.js"></script>
        -->
    <script src="/code/highstock.js"></script>
    <script src="/code/modules/exporting.js"></script>
    <script src="/code/modules/export-data.js"></script>
    <script src="/code/modules/accessibility.js"></script>

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
        <div class="form-group form-material">

            <table class="table table-hover table-striped">
                <thead>
                <tr>
                    <th>Parameter</th>
                    <th>Value</th>
                </tr>
                </thead>
                <tbody>

                <tr>
                    <td>SFP-SLA module:</td>
                    <td>
                        <?php echo app_module_label($moduleId); ?>
                    </td>
                </tr>

                </tbody>
            </table>

        </div>


        <figure class="highcharts-figure">
            <div id="container_back_to_back_delay" style="height: 1000px; min-width: 310px" >
                <img id="spiner_img" src="img/spiner.gif"  class="img-responsive center-block">

            </div>
        </figure>

        <script>
            Highcharts.setOptions({
                lang: {
                    loading: 'Loading...',
                    months: ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'],
                    weekdays: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
                    shortMonths: ['Jan', 'Feb', 'March', 'Apr', 'May', 'June', 'July', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
                    exportButtonTitle: "Export",
                    printButtonTitle: "Print",
                    rangeSelectorFrom: "From",
                    rangeSelectorTo: "To",
                    rangeSelectorZoom: "Period",
                    downloadPNG: 'Download PNG',
                    downloadJPEG: 'Download JPEG',
                    downloadPDF: 'Download PDF',
                    downloadSVG: 'Download SVG',
                    printChart: 'Print chart'
                }
            });
        </script>

        <script>

            Highcharts.getJSON('jsonl.php?id=<?php echo app_h($moduleId); ?>&n=600',    //3600
                function (data_init) {
                if (data_init==null)
                {
                    var elem = document.getElementById("spiner_img");
                    elem.remove()
                    var el = document.getElementById('container_back_to_back_delay');
                    var parentEl = el.parentNode;
                    newEl = document.createElement('h2');
                    attr = document.createAttribute("align");
                    attr.value = "center";    // set the created attribute value
                    newEl.setAttributeNode( attr );
                    newEl.innerHTML = 'No data to display';
                    parentEl.insertBefore(newEl, el);
                    return null

                }


                    if (data_init.length > 1) {
                        data_init.sort(function (a, b) {
                            return (a[0] - b[0]);
                        });
                        // force a redraw
                        //   data_init.update();
                    }

                    var load_to_lazer = [],
                        load_to_comm = [],
                        i;
                    for (i = 0; i < data_init.length; i += 1) {
                        load_to_lazer.push([data_init[i][0] * 1000, data_init[i][1]]);
                        load_to_comm.push([data_init[i][0] * 1000, data_init[i][2]]);
                    }


                    var this_chart = Highcharts.stockChart('container_back_to_back_delay', {
                        chart: {
                            events: {
                                load: function () {
                                    var series_1 = this.series[0],
                                        series_2 = this.series[1];

                                    setInterval(function () {
                                        Highcharts.getJSON('jsonl.php?id=<?php echo app_h($moduleId); ?>&n=1',
                                            function (points) {
                                                var i;
                                                for (i = 0; i < points.length; i += 1) {
                                                    series_1.addPoint([points[i][0] * 1000, points[i][1]], false, true, false);
                                                    series_2.addPoint([points[i][0] * 1000, points[i][2]], false, true, false);
                                                }
                                            })
                                        this_chart.redraw();
                                    }, 1000);
                                }
                            }
                        },
                        time: {
                            useUTC: false
                        },
                        rangeSelector: {
                            buttons: [{
                                count: 1,
                                type: 'minute',
                                text: '1M'
                            }, {
                                type: 'all',
                                text: 'All'
                            }]
                        },
                        xAxis: {
                            type: 'datetime',
                            tickPixelInterval: 150
                        },

                        yAxis: [ {
                                opposite: false,
                                labels: {
                                    align: 'right',
                                    x: -3
                                },
                                title: {
                                    text: 'Channel load from switch to laser (bit/s)'
                                },
                                min: 0,
                            offset: 0,
                                height: '45%',
                                lineWidth: 1
                            },
                            {
                                opposite: false,
                                labels: {
                                    align: 'right',
                                    x: -3
                                },
                                title: {
                                    text: 'Channel load from laser to switch (bit/s)'
                                },
                                min: 0,
                                top: '50%',
                                height: '45%',
                                offset: 0,
                                lineWidth: 1
                            }

                        ],
                        title: {
                            text: 'SFP-SLA module channel load'
                        },
                        exporting: {
                            enabled: false
                        },
                        tooltip: {
                            pointFormat: 'Channel load: <b>{point.y:.1f} bit/s</b>'
                        },
                        credits: {
                            enabled: false
                        },
                        series: [{
                            type: 'spline',
                            name: ',bit/s',
                            tooltip: {
                                valueDecimals: 2
                            },
                                data: load_to_lazer,
                                yAxis: 0
                        }, {
                            type: 'spline',
                            name: 'bit/s',
                            data: load_to_comm,
                            yAxis: 1
                        }]
                    });
                })
        </script>


    </div>

</div>

<!-- End Page -->
<!-- Footer -->
<footer class="site-footer">
    <div class="site-footer-legal">© 2019 <a href="http://Server SFP SLA.ru/" target="_blank">Server SFP SLA</a></div>

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
