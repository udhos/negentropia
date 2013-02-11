{{ define "title" }}Negentropia{{ end }}
{{ define "script" }}

{{if .Account}}
<script type="text/javascript" src="/js/ws.js"></script>
<script type="text/javascript" src="/js/lib/cookies.js"></script>

<script type="text/javascript">

var sid;

function start() {
	var prefix = "negentropia home javascript start(): ";
	console.log(prefix + "starting");

	sid = docCookies.getItem("sid");
	console.log(prefix + "cookie: sid=" + sid);

	statusElem = document.getElementById("ws_status");
	
	initWebSocket(statusElem);
	console.log(prefix + "done");
}

window.addEventListener("load", start, false);

</script>
{{end}}

{{ end }}
{{ define "content" }}

<div class="centerbox">

<h1>Welcome to Negentropia</h1>
<h3>This is Negentropia home location</h3>

<div id="ws_status"></div>

</div>


{{ end }}
