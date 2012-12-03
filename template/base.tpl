<!DOCTYPE html>
<html lang="en">
    <head>
        <title>{{ template "title" . }}</title>
    </head>
    <body>
        <section id="contents">
            {{ template "content" . }}
        </section>
        <footer id="footer">
            Footer goes here
        </footer>
    </body>
</html>