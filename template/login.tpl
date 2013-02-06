{{ define "title" }}Negentropia{{ end }}
{{ define "content" }}

<div class="centerbox">

<h1>Login</h1>

<form name="signup" action="{{.SignupPath}}" method="POST">
New user? <input type="submit" name="SignupButton" value="Signup">
<div><input type="hidden" name="Email"></div>
</form>

<form name="login" action="{{.LoginAuthPath}}" method="POST">

<div>Email: <input type="email" class="emailInput" spellcheck="false" name="Email" value="{{.EmailValue}}" placeholder="me@example.com" onchange="document.signup.Email.value = this.value;document.recover.Email.value = this.value;"></div>

<div>Password: <input type="password" name="Passwd"><input type="submit" name="LoginButton" value="Login"><span class="failmsg">{{.PasswdBadAuth}}</span></div>

<div><input type="submit" name="GoogleButton" value="Google Login"><span class="failmsg">{{.GoogleAuthMsg}}</span></div>

<div><input type="submit" name="FacebookButton" value="Facebook Login"><span class="failmsg">{{.FacebookAuthMsg}}</span></div>

<div><input type="submit" name="BrokenButton" value="Invalid Login"></div>
</form>

<form name="recover" action="{{.ResetPassPath}}" method="POST">
Forgot the password? <input type="submit" name="ResetPassButton" value="Password Recovery">
<div><input type="hidden" name="Email"></div>
</form>

</div>

{{ end }}
