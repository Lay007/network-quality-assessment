<?php
require_once 'app.php';
?>

<!DOCTYPE html>
<html class="no-js css-menubar" lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=0, minimal-ui">
    <meta name="description" content="SFP-SLA management console">
    <meta name="author" content="">
    <title>SFP-SLA Module Console</title>
    <link rel="apple-touch-icon" href="img/logo.svg">

    <link rel="shortcut icon" href="img/logo.svg" type="image/png">

    <!-- Stylesheets -->
    <link rel="stylesheet" href="css/bootstrap.min.css">
    <link rel="stylesheet" href="css/bootstrap-extend.min.css">
    <link rel="stylesheet" href="css/site.min.css">
    <!-- Plugins -->
    <link rel="stylesheet" href="vendor/animsition/animsition.css">
    <link rel="stylesheet" href="vendor/asscrollable/asScrollable.css">
    <link rel="stylesheet" href="vendor/waves/waves.css">
    <link rel="stylesheet" href="css/pages/login-v3.css">
    <!-- Fonts -->
    <link rel="stylesheet" href="fonts/material-design/material-design.min.css">
    <!-- Scripts -->
    <script src="vendor/modernizr/modernizr.js"></script>
    <script src="vendor/breakpoints/breakpoints.js"></script>
    <script>
        Breakpoints();
    </script>
</head>
<body class="page-login-v3 layout-full">

<!-- Page -->
<div class="page animsition vertical-align text-center" data-animsition-in="fade-in"
     data-animsition-out="fade-out">
    <div class="page-content vertical-align-middle">
        <div class="panel">
            <div class="panel-body">
                <div class="brand">
                    <a href="start.php"> <img class="brand-img" src="img/logo.svg" alt="Server SFP SLA">
                        <h2 class="brand-text font-size-18">SFP modules with SLA support</h2>
                    </a>
                </div>
                <form method="post" action="notepad.php" autocomplete="off">
                    <div class="form-group form-material floating">
                        <input type="text" class="form-control" name="log"/>
                        <label class="floating-label">Username</label>
                    </div>
                    <div class="form-group form-material floating">
                        <input type="password" class="form-control" name="pass"/>
                        <label class="floating-label">Password</label>
                    </div>
                    <button type="submit" class="btn btn-primary btn-block btn-lg margin-top-40">Sign in</button>
                </form>
                <p>Contact the system administrator to access this resource.</p>
            </div>
        </div>
        <footer class="page-copyright page-copyright-inverse">
            <p></p>
            <p>SFP-SLA Management Console.</p>
        </footer>
    </div>
</div>
<!-- End Page -->
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
<!-- Scripts -->
<script src="js/core.js"></script>
<script src="js/site.js"></script>
<script src="js/sections/menu.js"></script>
<script src="js/sections/menubar.js"></script>
<script src="js/configs/config-colors.js"></script>
<script src="js/components/asscrollable.js"></script>
<script src="js/components/animsition.js"></script>
<script src="js/components/material.js"></script>
<script>
    (function (document, window, $) {
        'use strict';
        var Site = window.Site;
        $(document).ready(function () {
            Site.run();
        });
    })(document, window, jQuery);
</script>
</body>
</html>
