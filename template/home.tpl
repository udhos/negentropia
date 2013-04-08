{{ define "title" }}Negentropia{{ end }}
{{ define "script" }}

{{if .Account}}
<script type="text/javascript" src="/js/lib/cookies.js"></script>
<script type="text/javascript" src="/js/lib/webgl-utils.js"></script>
<script type="text/javascript" src="/js/lib/webgl-debug.js"></script>
<script type="text/javascript" src="/js/lib/game-shim.js"></script>
<script type="text/javascript" src="/js/lib/Stats.js"></script>
<script type="text/javascript" src="/js/lib/gl-matrix-1.3.7.min.js"></script>

<script type="text/javascript" src="/js/main.js"></script>
<script type="text/javascript" src="/js/ws.js"></script>
<script type="text/javascript" src="/js/fetch.js"></script>
<script type="text/javascript" src="/js/shader.js"></script>
<script type="text/javascript" src="/js/buffer.js"></script>
<script type="text/javascript" src="/js/lost_context.js"></script>
<script type="text/javascript" src="/js/program.js"></script>
<script type="text/javascript" src="/js/skybox.js"></script>
<script type="text/javascript" src="/js/cubemap.js"></script>

<script type="text/javascript">

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

<span hidden id="wsUri">{{.Websocket}}</span>

<div class="centerbox">
<h1>Welcome to Negentropia</h1>
<h3>This is the home location</h3>
<div id="ws_status"></div>
</div>

<div id="canvasbox"></div>
<div id="framerate"></div>
<div id="control"></div>

<div id="toggle">switch to <a href="{{.HomeDartPath}}">dart</a></div>

{{ end }}
