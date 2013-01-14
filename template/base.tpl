<!DOCTYPE html>
<html lang="en">
    <head>
        <title>{{ template "title" . }}</title>
    </head>
    <body>
		<nav>
			{{if .Account}}{{.Account}} - <a href="/n/logout">logout</a>{{else}}<a href="/n/login">login</a>{{end}}
		</nav>
		
        <section id="contents">
            {{ template "content" . }}
        </section>
		
        <footer id="footer">
			Copyright (c) 2012 Negentropia
        </footer>
    </body>
</html>