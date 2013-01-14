{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}
<h1>Welcome to Negentropia</h1>

<form action="/n/loginAuth" method="POST">

<div>Email: <input type="email" spellcheck="false" name="Email" value="" placeholder="me@example.com"></div>

<div>Password: <input type="password" name="Passwd"><input type="submit" name="LoginButton" value="Login"><font color="red">{{.PasswdBadAuth}}</font></div>

<div><input type="submit" name="GoogleButton" value="Google Login"><font color="red">{{.GoogleAuthMsg}}</font></div>

<div><input type="submit" name="FacebookButton" value="Facebook Login"></div>

<div><input type="submit" name="BrokenButton" value="Invalid Login"></div>
</form>
{{ end }}
