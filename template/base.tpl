<!DOCTYPE html>
<html lang="en">
    <head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">	
		<link rel="stylesheet" href="/negentropia.css">
        <title>{{ template "title" . }}</title>
    </head>
    <body>
		<div class="bar">
		<span class="menu">
			{{if .ShowNavAccount}}
				{{if .Account}}<span id="account">{{.Account}}</span> - <a href="{{.LogoutPath}}">logout</a>{{else}}<a href="{{.LoginPath}}">login</a>{{end}}
			{{end}}
			{{if .ShowNavHome}}- <a href="{{.HomePath}}">home</a>{{end}}
			{{if .ShowNavSignup}}- <a href="{{.SignupPath}}">sign up</a>{{end}}
			{{if .ShowNavLogin}}- <a href="{{.LoginPath}}">login</a>{{end}}
			{{if .ShowNavLogout}}- <a href="{{.LogoutPath}}">logout</a>{{else}}{{end}}
		</span>
		</div>
		
        <section id="contents">
            {{ template "content" . }}
        </section>
		
        <footer id="footer">
			Copyright (c) 2013 Negentropia Team
        </footer>
		
		{{ template "script" . }}
    </body>
</html>