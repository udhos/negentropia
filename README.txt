
GENERAL BUILDING INSTRUCTIONS:

# 1. Clone the git repository:
git clone https://code.google.com/p/negentropia/

# Cloning will create the following tree:
negentropia/webserv
negentropia/webserv/bin
negentropia/webserv/pkg
negentropia/webserv/src
negentropia/webserv/src/negentropia
negentropia/webserv/src/negentropia/genid
negentropia/webserv/src/negentropia/webserv
negentropia/webserv/src/negentropia/webserv/handler

# 2. Set GOPATH to negentropia\webserv

# 3. Build and install (to negentropia\webserv\bin)
go install negentropia\webserv

# 4. Run
negentropia\webserv\bin\webserv.exe


BUILDING UNDER WINDOWS:

# windows git bash:

mkdir -p /c/tmp/devel

cd /c/tmp/devel

git clone https://code.google.com/p/negentropia/

# windows dos prompt:

set DEVEL=C:\tmp\devel

set GOPATH=%DEVEL%\negentropia\webserv

go install negentropia\webserv

%DEVEL%\negentropia\webserv\bin\webserv.exe

--THE END--
