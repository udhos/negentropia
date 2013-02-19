{{ define "title" }}Negentropia{{ end }}
{{ define "script" }}

{{if .Account}}
    <script type="application/dart" src="/dart/negentropia_home.dart"></script>
    <script src="/dart/dart.js"></script>
{{end}}
	
{{ end }}

{{ define "content" }}

<div class="centerbox">
<h1>Welcome to Negentropia</h1>
<h3>This is the DART home location</h3>
<div id="ws_status"></div>
</div>

<div id="canvasbox">
</div>

<a href="{{.HomePath}}">javascript</a>

{{ end }}
