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
<div id="control"></div>

<div id="toggle">switch to <a href="{{.HomeJSPath}}">javascript</a></div>

<div>mouse left click: select single item</div>
<div>hold shift + mouse left click: add/remove item to/from group selection</div>
<div>hold ctrl + drag mouse: band select multiple items</div>

{{ end }}
