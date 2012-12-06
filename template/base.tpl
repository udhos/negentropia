<!DOCTYPE html>
<html lang="en">
    <head>
        <title>{{ template "title" . }}</title>
    </head>
    <body>
		<nav>
			{{.Account}}
		</nav>
		
        <section id="contents">
            {{ template "content" . }}
        </section>
		
        <footer id="footer">
			Copyright (c) 2012 Negentropia
        </footer>
    </body>
</html>