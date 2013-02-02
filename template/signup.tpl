{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}
<h1>Sign up</h1>

<form action="{{.SignupProcessPath}}" method="POST">

<div>Name: <input type="text" name="Name" value="" placeholder="Your Name Here"></div>

<div>Email: <input type="email" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com"><span class="failmsg">{{.BadEmailMsg}}</span></div>

<div>Password: <input type="password" name="Passwd"><span class="failmsg">{{.BadPasswdMsg}}</span></div>

<div>Confirm: <input type="password" name="Confirm"><span class="failmsg">{{.BadConfirmMsg}}</span></div>

<div><input type="submit" name="SignupButton" value="Signup"><span class="failmsg">{{.BadSignupMsg}}</span></div>

</form>

<div class="donemsg">{{.SignupDoneMsg}}</div>

<p>Once you have signed up, please <a href="{{.ConfirmPath}}">confirm your email address</a>.</p>
{{ end }}
