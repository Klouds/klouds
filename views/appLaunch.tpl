	<!-- start content -->
	<div id="content">
		<div class="post">
			
			<div class="entry">
				<div align='center'>
				<form action="../../appLaunch" method="post" id="form1">
				  Name of App: {{str2html .FileSelect}}
				  Number to Launch: <input type="text" name="launchNumber"><br>
				</form>
				</div><div align='center'>
					
					<button type="submit" form="form1" value="Submit">Submit</button>
			</div></div>
			{{.Launching}}
		</div>
	</div>
	<!-- end content -->
