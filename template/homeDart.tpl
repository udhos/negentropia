{{ define "title" }}Negentropia{{ end }}

{{ define "script" }}

{{if .Account}}
    <script type="application/dart" src="/dart/negentropia_home.dart"></script>
    <script src="/dart/packages/browser/dart.js"></script>
{{end}}
	
{{ end }}

{{ define "content" }}

<span hidden id="wsUri">{{.Websocket}}</span>

<div class="centerbox">
<h1>Welcome to Negentropia</h1>
<h3>This is the DART home location</h3>
<div id="ws_status"></div>
</div>

<div id="canvasbox"></div>

<div id="framerate"></div>

<div id="toggle">return to <a href="{{.HomePath}}">javascript</a></div>

{{ end }}
