{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}
<h1>Confirm Email Adress</h1>

<form action="{{.ConfirmProcessPath}}" method="POST">

<div>Email: <input type="email" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com"><span class="failmsg">{{.BadEmailMsg}}</span></div>

<div>Confirmation id: <input type="text" name="ConfirmId"><span class="failmsg">{{.BadConfirmMsg}}</span></div>

<div><input type="submit" name="ConfirmButton" value="Confirm"></div>

</form>

<div class="donemsg">{{.ConfirmDoneMsg}}</div>
{{ end }}
