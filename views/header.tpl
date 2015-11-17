<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="description" content="">
	<meta name="keywords" content="Clouds Store">
	<meta name="author" content="Klouds.org">
	<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">

	<!-- SITE TITLE -->
	<title>Klouds</title>

	<!-- =========================
          FAV AND TOUCH ICONS
    ============================== -->
	<link rel="icon" href="/static/images/favicon.ico">
	<link rel="apple-touch-icon" href="/static/images/apple-touch-icon.png">
	<link rel="apple-touch-icon" sizes="72x72" href="/static/images/apple-touch-icon-72x72.png">
	<link rel="apple-touch-icon" sizes="114x114" href="/static/images/apple-touch-icon-114x114.png">

	<!-- =========================
         STYLESHEETS
    ============================== -->

	<!-- WEB FONTS -->
	<link href='http://fonts.googleapis.com/css?family=Roboto:100,300,100italic,400,300italic' rel='stylesheet'
		  type='text/css'>
	<link href='https://fonts.googleapis.com/css?family=Aclonica' rel='stylesheet' type='text/css'>
	<!-- CUSTOM STYLESHEETS -->
	<link rel="stylesheet" href="/static/css/styles.css">
	<link rel="stylesheet" href="/static/css/style-user.css">
	<!-- DEFAULT COLOR/ CURRENTLY USING -->

	<!-- RESPONSIVE FIXES -->
	<link rel="stylesheet" href="/static/css/responsive.css">
	<script src="/static/js/modernizr.js"></script>
	<!--[if lte IE 7]>
	<script src="/static/fonts/elegant-icons/lte-ie7.js"></script>
	<![endif]-->
	<!--[if lt IE 9]>
	<script src="/static/js/html5shiv.js"></script>
	<script src="/static/js/respond.min.js"></script>
	<![endif]-->


</head>

<body>
<!-- =========================
     PRE LOADER
============================== -->
<div class="preloader">
	<div class="status">&nbsp;</div>
</div>

<!-- =========================
     HEADER
============================== -->
<header class="header" data-stellar-background-ratio="0.5" id="home">

	<!-- COLOR OVER IMAGE -->
	<div class="main-content">
		<!-- To make header full screen. Use .full-screen class with color overlay. Example: <div class="color-overlay full-screen">  -->

		{{.Menu}}


		<!-- CONTAINER -->
		<div class="container">

			<!-- ONLY LOGO ON HEADER -->
			<div class="only-logo">
				<div class="navbar">
					<div class="navbar-header">
						<img src="/static/images/logo-color.png" alt="">
						<p class="font-common">speed up everything</p>
					</div>
				</div>
			</div>
			<!-- /END ONLY LOGO ON HEADER -->


			<div class="row">
				<div class="col-md-8 col-md-offset-2">

					<!-- HEADING AND BUTTONS -->
					<div class="intro-section">

						<!-- WELCOM MESSAGE -->
						<h1 class="intro">Put your ideas in a way with Klouds.</h1>
						<h5>Web Front End for interacting with Cloud-based REST APIs</h5>

						<!-- BUTTON -->
						<div class="buttons" id="download-button">

							<a href="#download" class="btn btn-default btn-lg standard-button"><i
										class="icon-app-download"></i>Sign Up Free</a>

						</div>
						<!-- /END BUTTONS -->

					</div>
					<!-- /END HEADNING AND BUTTONS -->

				</div>
			</div>
			<!-- /END ROW -->

		</div>
		<!-- /END CONTAINER -->
	</div>
	<!-- /END COLOR OVERLAY -->
</header>
<!-- /END HEADER -->
