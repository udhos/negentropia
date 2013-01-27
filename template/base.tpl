<!DOCTYPE html>
<html lang="en">
    <head>
        <title>{{ template "title" . }}</title>
    </head>
    <body>
		<nav>
			{{if .ShowNavAccount}}
			{{if .Account}}{{.Account}} - <a href="{{.LogoutPath}}">logout</a>{{else}}<a href="{{.LoginPath}}">login</a>{{end}}
			{{end}}
			{{if .ShowNavHome}}- <a href="{{.HomePath}}">home</a>{{end}}						
			{{if .ShowNavLogin}}- <a href="{{.LoginPath}}">login</a>{{end}}
			{{if .ShowNavLogout}}- <a href="{{.LogoutPath}}">logout</a>{{else}}{{end}}
		</nav>
		
        <section id="contents">
            {{ template "content" . }}
        </section>
		
        <footer id="footer">
			Copyright (c) 2013 Negentropia Team
        </footer>
    </body>
</html>