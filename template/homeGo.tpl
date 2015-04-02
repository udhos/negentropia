{{ define "title" }}Negentropia{{ end }}

{{ define "script" }}

{{if .Account}}
    <script async src="/negoc/negoc.js"></script>
{{end}}
	
{{ end }}

{{ define "content" }}

<span hidden id="wsUri">{{.Websocket}}</span>

<div class="centerbox">
<h1>Welcome to Negentropia</h1>
<h3>This is the GO home location</h3>
<div id="ws_status"></div>
</div>

<div id="canvasbox"></div>
<div id="debug"></div>
<div id="framerate"></div>
<div id="control"></div>

<div id="toggle">switch to <a href="{{.HomePath}}">dart</a></div>
<div id="toggle">switch to <a href="{{.HomeJSPath}}">javascript</a></div>

<div>mouse left click: select single item</div>
<div>hold shift + mouse left click: add/remove item to/from group selection</div>
<div>hold ctrl + drag mouse: band select multiple items</div>
<div>t: toggle camera tracking on/off</div>
<div>hold mouse right button + drag mouse: camera rotate</div>
<div>mouse wheel: camera zoom in/out</div>
<div>space: restore camera default orientation</div>
<div>F2: switch selected units' mission</div>

{{ end }}
