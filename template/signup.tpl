{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}
<h1>Sign up</h1>

<form action="{{.SignupProcessPath}}" method="POST">

<div>Name: <input type="text" name="Name" value="" placeholder="Your Name Here"></div>

<div>Email: <input type="email" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com"></div>

<div>Password: <input type="password" name="Passwd"><font color="red">{{.BadPasswdMsg}}</font></div>

<div>Confirm: <input type="password" name="Confirm"><font color="red">{{.BadConfirmMsg}}</font></div>

<div><input type="submit" name="SignupButton" value="Signup"><font color="red">{{.BadSignupMsg}}</font></div>

</form>
{{ end }}
