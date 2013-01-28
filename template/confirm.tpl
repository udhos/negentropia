{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}
<h1>Sign up</h1>

<form action="{{.ConfirmProcessPath}}" method="POST">

<div>Email: <input type="email" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com"><font color="red">{{.BadEmailMsg}}</font></div>

<div>Confirmation id: <input type="text" name="ConfirmId"><font color="red">{{.BadConfirmMsg}}</font></div>

<div><input type="submit" name="ConfirmButton" value="Confirm"></div>

</form>

<div><font color="blue">{{.ConfirmDoneMsg}}</font></div>
{{ end }}
