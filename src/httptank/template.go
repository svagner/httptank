package main

const (
	TANK_TMPL = `
<html>
<head>
<title>TANK - {{ .Title }}</title>
  <!-- Latest compiled and minified CSS -->
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
  <!-- Optional theme -->
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap-theme.min.css">
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
  <!-- Latest compiled and minified JavaScript -->
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js"></script>
  <script>
    $(document).on('click', '.panel-heading button', function(e) {
      var $this = $(this);
      var icon = $this.find('i');
      if (icon.hasClass('glyphicon-plus')) {
        $this.find('i').removeClass('glyphicon-plus').addClass('glyphicon-minus');
      } else {
        $this.find('i').removeClass('glyphicon-minus').addClass('glyphicon-plus');
      }
    });
  </script>
  <style>
    .clickable {
      cursor: pointer;
    }
    .clickable .glyphicon {
      background: rgba(0, 0, 0, 0.15);
      display: inline-block;
      padding: 6px 12px;
      border-radius: 4px
    }
    .panel-heading span {
      margin-top: -23px;
      font-size: 15px;
      margin-right: -9px;
    }
    .panel-heading button  {
      margin-top: -25px;
    }
  </style>
<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.selection.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.resize.min.js"></script>
<script type="text/javascript">

$(document).on('submit','#settings',function(event){    
  event.preventDefault()
  var formData = $(this).serialize();
  console.log(formData);
  $.ajax({
    url: "/settings",
    data: formData,
    datatype: "json",
    type: "POST",
    success: function(data) {}
  });
});

(function() {
	var data = [
		{ label: "count", data: {{ .Error }} },
		{ label: "error", data: {{ .Count }} },
	];

	var dataTime = [
		{ label: "Avg", data: {{ .AvgTime }} },
		{ label: "Max", data: {{ .MaxTime }} },
		{ label: "Min", data: {{ .MinTime }} },
	];

	var dataErrors = [
		{ label: "50x", data: {{ .Error50x }} },
		{ label: "40x", data: {{ .Error40x }} },
		{ label: "Timeout", data: {{ .ErrorTimeout }} },
		{ label: "Other", data: {{ .ErrorOther }} },
	];

	var options = {
		legend: {
			position: "nw",
			noColumns: 2,
			backgroundOpacity: 0.2
		},
		yaxis: {
			tickFormatter: function(val) { return val; }
		},
		xaxis: {
			tickFormatter: function(val) { return val + "s"; }
		},
		selection: {
			mode: "x"
		},
	};

	var optionsTime = {
		legend: {
			position: "nw",
			noColumns: 2,
			backgroundOpacity: 0.2
		},
		yaxis: {
			tickFormatter: function(val) { return val; }
		},
		xaxis: {
			tickFormatter: function(val) { return val + "s"; }
		},
		selection: {
			mode: "x"
		},
		lines: {
			show: true,
			fill: true,
			lineWidth: 1
		},
		points: {
		  show: true
		},
		grid: {
		  hoverable: true,
		  clickable: true
		},
	};



	$(document).ready(function() {
		var plot = $.plot("#placeholder", data, options);
		var plotTime = $.plot("#placeholderTime", dataTime, optionsTime);
		var plotErrors = $.plot("#placeholderErrors", dataErrors, optionsTime);

		var overview = $.plot("#overview", data, {
			legend: { show: true },
			series: {
				lines: {
					show: true,
					lineWidth: 1
				},
				shadowSize: 0
			},
			grid: {
				  hoverable: true,
					  clickable: true
			},
			xaxis: {
				ticks: [],
				min: 0,
				autoscaleMargin: 0.1
			},
			yaxis: {
				ticks: [],
				min: 0,
				autoscaleMargin: 0.1
			},
			selection: {
				mode: "x"
			}
		});

		$("<div id='tooltip'></div>").css({
		  position: "absolute",
		  display: "none",
		  border: "1px solid #fdd",
		  padding: "2px",
		  "background-color": "#fee",
		  opacity: 0.80
		}).appendTo("body");



		// now connect the two
		$("#placeholder").bind("plotselected", function (event, ranges) {

			// do the zooming
			$.each(plot.getXAxes(), function(_, axis) {
				var opts = axis.options;
				opts.min = ranges.xaxis.from;
				opts.max = ranges.xaxis.to;
			});
			plot.setupGrid();
			plot.draw();
			plot.clearSelection();

			// don't fire event on the overview to prevent eternal loop

			overview.setSelection(ranges, true);
			plotTime.setSelection(ranges, true);
			plotErrors.setSelection(ranges, true);
		});

		$("#placeholderErrors").bind("plothover", function (event, pos, item) {
			if (item) {
			  var x = item.datapoint[0].toFixed(2),
			  y = item.datapoint[1].toFixed(2);

			  $("#tooltip").html(item.series.label + " of " + x + "s. = " + y)
			  .css({top: item.pageY+5, left: item.pageX+5})
			  .fadeIn(200);
			} else {
			  $("#tooltip").hide();
			}
		});

		$("#placeholderTime").bind("plothover", function (event, pos, item) {
			if (item) {
			  var x = item.datapoint[0].toFixed(2),
			  y = item.datapoint[1].toFixed(2);

			  $("#tooltip").html(item.series.label + " of " + x + "s. = " + y + "ms.")
			  .css({top: item.pageY+5, left: item.pageX+5})
			  .fadeIn(200);
			} else {
			  $("#tooltip").hide();
			}
		});

		$("#placeholderErrors").bind("plotselected", function (event, ranges) {
			// do the zooming
			$.each(plotErrors.getXAxes(), function(_, axis) {
				var opts = axis.options;
				opts.min = ranges.xaxis.from;
				opts.max = ranges.xaxis.to;
			});
			plotErrors.setupGrid();
			plotErrors.draw();
			plotErrors.clearSelection();

			// don't fire event on the overview to prevent eternal loop

			overview.setSelection(ranges, true);
			plot.setSelection(ranges, true);
			plotTime.setSelection(ranges, true);
		});

		$("#placeholderTime").bind("plotselected", function (event, ranges) {
			// do the zooming
			$.each(plotTime.getXAxes(), function(_, axis) {
				var opts = axis.options;
				opts.min = ranges.xaxis.from;
				opts.max = ranges.xaxis.to;
			});
			plotTime.setupGrid();
			plotTime.draw();
			plotTime.clearSelection();

			// don't fire event on the overview to prevent eternal loop

			overview.setSelection(ranges, true);
			plot.setSelection(ranges, true);
			plotErrors.setSelection(ranges, true);
		});

		$("#overview").bind("plotselected", function (event, ranges) {
			plot.setSelection(ranges);
			plotTime.setSelection(ranges);
			plotErrors.setSelection(ranges);
		});

		// refresh data every second
		pullAndRedraw();

		function pullAndRedraw() {
			$.get(window.location.href + 'graph.json', function(graphData) {
				var data = [
					{ label: "count", data: graphData.Count },
					{ label: "error", data: graphData.Error }
				];
				var dataTime = [
					{ label: "Avg", data: graphData.AvgTime },
					{ label: "Max", data: graphData.MaxTime },
					{ label: "Min", data: graphData.MinTime }
				];
				var dataErrors = [
					{ label: "50x", data: graphData.Error50x },
					{ label: "40x", data: graphData.Error40x },
					{ label: "Timeout", data: graphData.ErrorTimeout },
					{ label: "Other", data: graphData.ErrorOther }
				];

				plot.setData(data);
				plot.setupGrid();
				plot.draw();

				plotTime.setData(dataTime);
				plotTime.setupGrid();
				plotTime.draw();

				plotErrors.setData(dataErrors);
				plotErrors.setupGrid();
				plotErrors.draw();

				overview.setData(data);
				overview.setupGrid();
				overview.draw();

				setTimeout(pullAndRedraw, 1000);
			})
		}
	});
})();
</script>
<style>
#content {
	margin: 0 auto;
	padding: 10px;
}

#export {
	float: right;
}

.tab-pane {
	box-sizing: border-box;
	width: 1200px;
	height: 450px;
	//padding: 20px 15px 15px 15px;
	margin: 15px auto 30px auto;
	border: 1px solid #ddd;
	background: #fff;
	background: linear-gradient(#f6f6f6 0, #fff 50px);
	background: -o-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -ms-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -moz-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -webkit-linear-gradient(#f6f6f6 0, #fff 50px);
	box-shadow: 0 3px 10px rgba(0,0,0,0.15);
	-o-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-ms-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-moz-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-webkit-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
}
.demo-container {
	box-sizing: border-box;
	width: 1200px;
	height: 450px;
	padding: 20px 15px 15px 15px;
	margin-left: auto;
	margin-right: auto;
	//margin: 15px auto 30px auto;
	//border: 1px solid #ddd;
	background: #fff;
	background: linear-gradient(#f6f6f6 0, #fff 50px);
	background: -o-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -ms-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -moz-linear-gradient(#f6f6f6 0, #fff 50px);
	background: -webkit-linear-gradient(#f6f6f6 0, #fff 50px);
	box-shadow: 0 3px 10px rgba(0,0,0,0.15);
	-o-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-ms-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-moz-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
	-webkit-box-shadow: 0 3px 10px rgba(0,0,0,0.1);
}

.demo-placeholder {
	width: 100%;
	height: 100%;
	font-size: 14px;
	line-height: 1.2em;
}
</style>
</head>
<body>


<div class="container">
  <div class="row">
    <div class="col-md-12">
      <div class="panel panel-primary">
	<div class="panel-heading">
	  <h3 class="panel-title">Settings</h3>
	  <button class="btn btn-primary pull-right" type="button" data-toggle="collapse" data-target="#collapseExample" aria-expanded="false" aria-controls="collapseExample">
	    <i class="glyphicon glyphicon-plus"></i>
	  </button>
	</div>
	<div class="collapse" id="collapseExample">
	  <div class="panel-body">
	    <form id="settings" name="settings" class="form-horizontal container" role="form" method="post" action="">
	    <div class="row">
  	      <div class="form-group col-sm-6">
  	        <label for="url" class="col-sm-4 control-label">URL:</label>
  	        <div class="col-sm-8">
  	          <input type="text" class="form-control" id="url" name="url" placeholder="http://localhost" value="{{ .Settings.Url }}">
  	        </div>
  	      </div>
  	      <div class="form-group col-sm-6">
  	        <label for="time" class="col-sm-4 control-label">Test time (sec):</label>
  	        <div class="col-sm-8">
  	          <input type="number" class="form-control" id="time" name="time" placeholder="10" value="{{ if eq .Settings.Time 0 }}10{{else}}{{ .Settings.Time }}{{ end }}">
  	        </div>
  	      </div>
  	    </div>
	    <div class="row">
  	      <div class="form-group col-sm-6">
  	        <label for="parallel" class="col-sm-4 control-label">Parallel queries:</label>
  	        <div class="col-sm-8">
  	          <input type="number" class="form-control" id="parallel" name="parallel" placeholder="1" value="{{ if eq .Settings.Count 0 }}1{{else}}{{.Settings.Count}}{{end}}">
  	        </div>
  	      </div>
  	      <div class="form-group col-sm-6">
  	        <label for="timeout" class="col-sm-4 control-label">Query timeout (ms):</label>
  	        <div class="col-sm-8">
  	          <input type="number" class="form-control" id="timeout" name="timeout" placeholder="100" value="{{if eq .Settings.Timeout 0}}100{{else}}{{.Settings.Timeout}}{{end}}">
  	        </div>
  	      </div>
	    </div>
	    <div class="row">
  	      <div class="form-group col-sm-12">
	      <h4 align="center">Basic Authorization</h4>
  	      </div>
	    </div>
	    <div class="row">
  	      <div class="form-group col-sm-6">
  	        <label for="username" class="col-sm-4 control-label">Username:</label>
  	        <div class="col-sm-8">
  	          <input type="text" class="form-control" id="username" name="username" placeholder="user" value="{{ if eq .Settings.Username "" }}{{.Settings.Username}}{{end}}">
  	        </div>
  	      </div>
  	      <div class="form-group col-sm-6">
  	        <label for="password" class="col-sm-4 control-label">Password:</label>
  	        <div class="col-sm-8">
  	          <input type="password" class="form-control" id="password" name="password" placeholder="password" value="{{if ne .Settings.Password ""}}{{.Settings.Password}}{{end}}">
  	        </div>
  	      </div>
	    </div>
	    <div class="row">
  	      <div class="form-group col-sm-12">
	      <h4 align="center">HTTP Client Parameters</h4>
  	      </div>
	    </div>
	    <div class="row">
  	      <div class="form-group col-sm-6">
  	        <label for="useragent" class="col-sm-4 control-label">UserAgent:</label>
  	        <div class="col-sm-8">
  	          <input type="text" class="form-control" id="useragent" name="useragent" placeholder="Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36" value="{{ if ne .Settings.Useragent "" }}{{.Settings.Useragent}}{{else}}Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36{{end}}">
  	        </div>
  	      </div>
  	      <div class="form-group col-sm-6">
  	        <label for="cookie" class="col-sm-4 control-label">Cookie:</label>
  	        <div class="col-sm-8">
		<input type="text" class="form-control" id="cookie" name="cookie" placeholder="name1:value1:/:.google.com,name2:value2:/test:.google.net" value="{{if ne .Settings.Cookie ""}}{{.Settings.Cookie}}{{end}}">
  	        </div>
		<p class="help-block">name1:value1:/:.google.com - where '/' - path, .google.com - domain</p>
  	      </div>
	    </div>
	    <div class="row">
  	      <div class="form-group">
  	          <div class="col-sm-5 col-sm-offset-2">
  	            <input id="submit" name="submit" type="submit" value="Send" class="btn btn-primary">
		  </div>
  	      </div>
	    </div>
  	    </form>
	  </div>
	</div>
      </div>
    </div>
  </div>
</div>

<div id="export">
	<a href="/graph.json">json</a>
</div>
<div id="content">
  <ul class="nav nav-tabs">
    <li class="active"><a data-toggle="tab" href="#count">Count</a></li>
    <li><a data-toggle="tab" href="#timing">Timing</a></li>
    <li><a data-toggle="tab" href="#errors">Errors</a></li>
  </ul>
  <div class="tab-content">
    <div id="count" class="tab-pane fade in active">
	<div class="demo-container">
		<div id="placeholder" class="demo-placeholder"></div>
	</div>
    </div>
    <div id="timing" class="tab-pane">
	<div class="demo-container">
		<div id="placeholderTime" class="demo-placeholder"></div>
	</div>
    </div>
    <div id="errors" class="tab-pane">
	<div class="demo-container">
		<div id="placeholderErrors" class="demo-placeholder"></div>
	</div>
    </div>
  </div>

	<div class="demo-container" style="height:150px;">
		<div id="overview" class="demo-placeholder"></div>
	</div>

	<p>The smaller plot is linked to the main plot, so it acts as an overview. Try dragging a selection on either plot, and watch the behavior of the other.</p>

</div>

<pre><b>Legend</b>

count: queries count 
error: queries count with bad state 

Dinamic parameters for queries:
- .Param - string from file which set by 'file' argv
</pre>
</body>
</html>
	`
)
