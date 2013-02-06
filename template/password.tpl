{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}

<div class="centerbox">

<h1>Password Recovery</h1>

<form action="{{.ResetPassProcessPath}}" method="POST">

<div>Email: <input type="email" class="emailInput" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com" ><span class="failmsg">{{.BadEmailMsg}}</span></div>

<div><input type="submit" name="ResetPassButton" value="Reset Password"></div>

</form>

<div class="donemsg">{{.ResetPassDoneMsg}}</div>

<p>Once you have requested the password recovery code, please <a href="{{.ResetPassConfirmPath}}">enter the new password</a>.</p>

</div>

{{ end }}
