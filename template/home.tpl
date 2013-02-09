{{ define "title" }}Negentropia{{ end }}
{{ define "script" }}

{{if .Account}}
<script type="text/javascript" src="/js/ws.js"></script>
<script type="text/javascript" src="/js/lib/cookies.js"></script>

<script type="text/javascript">

var sid;

function start() {
	console.log("starting");

	sid = docCookies.getItem("sid");
	console.log("cookie: sid=" + sid);

	initWebSocket();
	console.log("done");
}

window.addEventListener("load", start, false);

</script>
{{end}}

{{ end }}
{{ define "content" }}

<div class="centerbox">

<h1>Welcome to Negentropia</h1>
<h3>This is Negentropia home location</h3>
</div>

{{ end }}
