{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}
<h1>Password Recovery</h1>

<form action="{{.ResetPassProcessPath}}" method="POST">

<div>Email: <input type="email" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com" ><font color="red">{{.BadEmailMsg}}</font></div>

<div><input type="submit" name="ResetPassButton" value="Reset Password"><font color="red"></font></div>

</form>

<div><font color="blue">{{.ResetPassDoneMsg}}</font></div>

<div>Once you have request password recovery, please <a href="{{.ResetPassConfirmPath}}">enter the new password</a>.</div>
{{ end }}
