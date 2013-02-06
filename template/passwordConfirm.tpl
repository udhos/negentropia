{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}

<div class="centerbox">

<h1>Confirm Password Recovery</h1>

<form action="{{.ResetPassConfirmProcessPath}}" method="POST">

<div>Email: <input type="email" class="emailInput" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com"><span class="failmsg">{{.BadEmailMsg}}</span></div>

<div>Confirmation id: <input type="text" name="ConfirmId" value="{{.ConfirmIdValue}}"><span class="failmsg">{{.BadConfirmIdMsg}}</span></div>

<div>New Password: <input type="password" name="Passwd"><span class="failmsg">{{.BadPasswdMsg}}</span></div>

<div>Confirm Password: <input type="password" name="Confirm"><span class="failmsg">{{.BadConfirmMsg}}</span></div>

<div><input type="submit" name="ResetPassConfirmButton" value="Change Password"></div>

</form>

<div class="donemsg">{{.ResetPassConfirmDoneMsg}}</div>

</div>

{{ end }}
