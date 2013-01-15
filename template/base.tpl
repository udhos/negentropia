<!DOCTYPE html>
<html lang="en">
    <head>
        <title>{{ template "title" . }}</title>
    </head>
    <body>
		<nav>
			{{if .ShowNavAccount}}
			{{if .Account}}{{.Account}} - <a href="/n/logout">logout</a>{{else}}<a href="/n/login">login</a>{{end}}
			{{end}}
			{{if .ShowNavHome}}- <a href="/n/">home</a>{{end}}						
			{{if .ShowNavLogin}}- <a href="/n/login">login</a>{{end}}
			{{if .ShowNavLogout}}- <a href="/n/logout">logout</a>{{else}}{{end}}
		</nav>
		
        <section id="contents">
            {{ template "content" . }}
        </section>
		
        <footer id="footer">
			Copyright (c) 2012 Negentropia
        </footer>
    </body>
</html>