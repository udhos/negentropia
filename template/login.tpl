{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}
<h1>Login</h1>

<form name="signup" action="{{.SignupPath}}" method="POST">
New user? <input type="submit" name="SignupButton" value="Signup">
<div><input type="hidden" name="Email"></div>
</form>

<form name="login" action="{{.LoginAuthPath}}" method="POST">

<div>Email: <input type="email" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com" onblur="document.signup.Email.value = this.value;document.recover.Email.value = this.value;"></div>

<div>Password: <input type="password" name="Passwd"><input type="submit" name="LoginButton" value="Login"><font color="red">{{.PasswdBadAuth}}</font></div>

<div><input type="submit" name="GoogleButton" value="Google Login"><font color="red">{{.GoogleAuthMsg}}</font></div>

<div><input type="submit" name="FacebookButton" value="Facebook Login"><font color="red">{{.FacebookAuthMsg}}</font></div>

<div><input type="submit" name="BrokenButton" value="Invalid Login"></div>
</form>

<form name="recover" action="{{.ResetPassPath}}" method="POST">
Forgot the password? <input type="submit" name="ResetPassButton" value="Password Recovery">
<div><input type="hidden" name="Email"></div>
</form>

{{ end }}
