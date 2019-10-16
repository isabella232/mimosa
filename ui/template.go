package ui

func getHTML() string {
	s := `<!DOCTYPE html>
	<html>

	<head>
	  <meta charset="utf-8">
	  <title>Host Data</title>
	  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.10/semantic.min.css">
	  <!-- <link rel="stylesheet" href="semantic/dist/semantic.min.css"> -->


	  <style type="text/css">
		body>.ui.container {
		  margin-top: 3em;
		}

		.ui.container>h1 {
		  font-size: 3em;
		  text-align: center;
		  font-weight: normal;
		}

		.ui.container>h2.dividing.header {
		  font-size: 2em;
		  font-weight: normal;
		  margin: 4em 0em 3em;
		}

		.ui.table {
		  table-layout: fixed;
		}
	  </style>

	</head>

	<body>
	  <div class="ui container">
		<h2 class="ui horizontal divider header">
		  <i class="globe icon"></i>
		  Host Data
		</h2>

		<div class="ui segment">

<!--
		  <form method="get" id="search"></form>
		  <form method="get" id="reset"></form>

		  <div class="ui horizontal segments">
			<div class="ui segment">
			  <div class="ui action input">
				<input form="search" type="text" name="filter" value="" placeholder="Search...">
				<button form="search" class="ui button" type="submit">Search</button>
				<button form="reset" class="ui button" type="submit">Reset</button>
			  </div>
			</div>
		  </div>
-->

		  <table class="ui celled striped table">
			<thead>
			  <tr>
				<th>Name</th>
				<th>PublicDNS</th>
				<th>PublicIP</th>
			  </tr>
			</thead>
			<tbody>
			  {{range .}}
			  <tr>
				<td>{{.Name}}</td>
				<td>{{.PublicDNS}}</td>
				<td>{{.PublicIP}}</td>
			  </tr>
			  {{end}}
			</tbody>
		  </table>
		</div>

	  </div>
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.10/semantic.min.js"></script>
	  <!-- <script src="semantic/dist/semantic.min.js"></script> -->
	</body>

	</html>`
	return s
}
