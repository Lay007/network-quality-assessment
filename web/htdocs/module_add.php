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


            <div class="row">
                <div class="col-lg-12">
                    <div class="panel panel-bordered">
                        <div class="panel-heading">
                            <h3 class="panel-title"><i class="icon md-link" aria-hidden="true"></i>Add module
                                SFP-SLA</h3>
                        </div>
                        <div class="panel-body">
                            <form method="post" action="action_add_module_go.php" class="form-horizontal"
                                 autocomplete="off">
                                <div class="form-group form-material">
                                    <label class="col-sm-3 control-label">SFP-SLA module name</label>
                                    <div class="col-sm-9">
                                        <input type="text" name=module_name class="form-control"/>
                                    </div>
                                </div>

                                <div class="form-group form-material">
                                    <label class="col-sm-3 control-label">SFP-SLA module IP address </label>
                                    <div class="col-sm-9">
                                        <input type="text" name="module_ip" pattern="^([0-9]{1,3}\.){3}[0-9]{1,3}$" class="form-control"/>
                                    </div>
                                </div>
                                <div class="form-group form-material">
                                    <label class="col-sm-3 control-label">Firmware version</label>
                                    <div class="col-sm-9">
                                        <input type="text" name="module_ver" class="form-control"/>
                                    </div>
                                </div>




                                <div class="form-group form-material">
                                    <label class="col-sm-3 control-label">Additional description</label>
                                    <div class="col-sm-9">
                                        <textarea class="form-control" name="module_desc" rows="5"></textarea>
                                    </div>
                                </div>

                                <div class="text-right">
                                    <button type="submit" class="btn btn-primary">Add module</button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>



        </div>
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

