{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}
<h1>Confirm Password Recovery</h1>

<form action="{{.ResetPassConfirmProcessPath}}" method="POST">

<div>Email: <input type="email" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com"><font color="red">{{.BadEmailMsg}}</font></div>

<div>Confirmation id: <input type="text" name="ConfirmId" value="{{.ConfirmIdValue}}"><font color="red">{{.BadConfirmIdMsg}}</font></div>

<div>New Password: <input type="password" name="Passwd"><font color="red">{{.BadPasswdMsg}}</font></div>

<div>Confirm Password: <input type="password" name="Confirm"><font color="red">{{.BadConfirmMsg}}</font></div>

<div><input type="submit" name="ResetPassConfirmButton" value="Change Password"></div>

</form>

<div><font color="blue">{{.ResetPassConfirmDoneMsg}}</font></div>
{{ end }}
