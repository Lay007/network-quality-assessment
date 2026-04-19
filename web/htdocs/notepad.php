<?php
require_once 'db.php';
?>

<!DOCTYPE html>
<html class="no-js css-menubar" lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=0, minimal-ui">
    <meta name="description" content="SFP-SLA management console">
    <meta name="author" content="">
    <title>Authentication error</title>
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
    <link rel="stylesheet" href="css/pages/errors.css">
    <!-- Fonts -->
    <link rel="stylesheet" href="fonts/material-design/material-design.min.css">
    <!-- Scripts -->
    <script src="vendor/modernizr/modernizr.js"></script>
    <script src="vendor/breakpoints/breakpoints.js"></script>
    <script>
        Breakpoints();
    </script>
</head>
<body class="page-error page-error-503 layout-full">

<!-- Page -->
<div class="page animsition vertical-align text-center" data-animsition-in="fade-in"
     data-animsition-out="fade-out">
    <div class="page-content vertical-align-middle">
        <header>
            <?php
            $login = app_post_string('log');
            $pass = app_post_string('pass');

            if ($login !== '' && $pass !== '') {
                $stmt = app_stmt('SELECT `password` FROM `users` WHERE `login` = ?', 's', [$login]);
                $search = mysqli_stmt_get_result($stmt);
                $row = $search ? mysqli_fetch_assoc($search) : null;

                if (!$row) {
                    echo "<h2 class='animation-slide-top'>Authentication error</h2><p>User not found</p>";
                } else {
                    $storedPassword = $row['password'];
                    $passwordValid = password_verify($pass, $storedPassword);

                    if (!$passwordValid) {
                        echo "<h2 class='animation-slide-top'>Authentication error</h2><p>Incorrect password</p>";
                    } else {
                        $_SESSION['user_in'] = $login;
                        if (!isset($_SESSION['time'])) {
                            $_SESSION['time'] = date('H:i:s');
                        }

                        echo "<meta http-equiv='Refresh' Content='0;URL=menu.php'>";
                    }
                }
            } else {
                echo "<h2 class='animation-slide-top'>Authentication error</h2><p>Login or password is missing</p>";
            }
            ?>

        </header>
        <p class="error-advise">Return to the login page</p>
        <a class="btn btn-primary btn-round" href="index.php">GO TO LOGIN PAGE</a>
        <footer class="page-copyright">
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
