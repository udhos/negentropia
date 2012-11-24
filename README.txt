
REQUIREMENTS

# Install Go Language
# http://golang.org/doc/install

# Install GIT - required to clone the source tree
# http://git-scm.com/downloads

# Install Mercurial - required for packages such as: go get code.google.com/p/goauth2/oauth
# http://mercurial.selenic.com/downloads/


GENERAL BUILDING INSTRUCTIONS:

# 1. Clone the git repository:
git clone https://code.google.com/p/negentropia/

# 2. Set GOPATH to negentropia\webserv

# 3. install goauth2
go get code.google.com/p/goauth2/oauth

# 4. Build and install (to negentropia\webserv\bin)
go install negentropia\webserv

# 5. Run
# Under Linux:
negentropia\webserv\bin\webserv
# Under Windows:
negentropia\webserv\bin\webserv.exe


BUILDING UNDER WINDOWS:

# windows git bash:

mkdir -p /c/tmp/devel

cd /c/tmp/devel

git clone https://code.google.com/p/negentropia/

# windows dos prompt:

set DEVEL=C:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

# install goauth2
go get code.google.com/p/goauth2/oauth

# build
go install negentropia\webserv

# run
%DEVEL%\negentropia\webserv\bin\webserv.exe

--THE END--
