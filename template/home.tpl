{{ define "title" }}Negentropia{{ end }}
{{ define "script" }}

{{if .Account}}
<script type="text/javascript" src="/js/main.js"></script>
<script type="text/javascript" src="/js/ws.js"></script>
<script type="text/javascript" src="/js/lib/cookies.js"></script>
<script type="text/javascript" src="/js/lib/webgl-utils.js"></script>
<script type="text/javascript" src="/js/lib/webgl-debug.js"></script>

<script type="text/javascript">

// Global variables
var neg = {
	debugLostContext: true,
	debugWebGL: true,
	wsUri: {{.Websocket}}
};

function start() {
	var prefix = "negentropia home javascript start(): ";
	console.log(prefix + "starting");
	
	main(); // main.js

	console.log(prefix + "done");
}

window.addEventListener("load", start, false);

</script>
{{end}}

{{ end }}

{{ define "content" }}

<div class="centerbox">
<h1>Welcome to Negentropia</h1>
<h3>This is the home location</h3>
<div id="ws_status"></div>
</div>

<div id="canvasbox">
<!--
<canvas id="main_canvas">
Browser missing &lt;canvas&gt; support!
</canvas>
-->
</div>

{{ end }}
