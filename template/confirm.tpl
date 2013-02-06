{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}

<div class="centerbox">

<h1>Confirm Email Address</h1>
<h3>Enable new account</h3>

<form action="{{.ConfirmProcessPath}}" method="POST">

<div>Email: <input type="email" class="emailInput" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com"><span class="failmsg">{{.BadEmailMsg}}</span></div>

<div>Confirmation id: <input type="text" name="ConfirmId"><span class="failmsg">{{.BadConfirmMsg}}</span></div>

<div><input type="submit" name="ConfirmButton" value="Confirm"></div>

</form>

<div class="donemsg">{{.ConfirmDoneMsg}}</div>

</div>

{{ end }}
