
REQUIREMENTS

# Install Go Language
# http://golang.org/doc/install

# Install GIT - required to clone the source tree
# http://git-scm.com/downloads

# Install Mercurial - required for packages such as: go get code.google.com/p/goauth2/oauth
# http://mercurial.selenic.com/downloads/

# Install Memcached
# http://memcached.org/
#
# memcached 1.4.2 for Windows 7 64-bit:
# http://www.urielkatz.com/archive/detail/memcached-64-bit-windows/
# http://www.urielkatz.com/projects/memcached-win64/memcached-win64.zip
# http://dl.dropbox.com/u/103946/memcached-win64.zip
# http://fajarmf.net/apps/memcached-win64.zip


GENERAL BUILDING INSTRUCTIONS:

# 1. Clone the git repository:
git clone https://code.google.com/p/negentropia/

# 2. Set GOPATH to negentropia\webserv

# 3. install goauth2
go get code.google.com/p/goauth2/oauth

# 4. Install memcache client library
go get github.com/bradfitz/gomemcache/memcache

# 5. Build and install (to negentropia\webserv\bin)
go install negentropia\webserv

# 6. Start memcached
memcached -vv -p 11211

# 7. Run
# Under Linux:
negentropia\webserv\bin\webserv
# Under Windows:
negentropia\webserv\bin\webserv.exe


BUILDING UNDER WINDOWS:

## windows git bash:

export DEVEL=/c/tmp/devel
export GOPATH=$DEVEL/negentropia/webserv

mkdir -p $DEVEL
cd $DEVEL

git clone https://code.google.com/p/negentropia/

go get github.com/bradfitz/gomemcache/memcache

## windows dos prompt:

# start memcached
memcached.exe -vv -p 11211

set DEVEL=C:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

# install goauth2
go get code.google.com/p/goauth2/oauth

# build
go install negentropia\webserv

# run
%DEVEL%\negentropia\webserv\bin\webserv.exe

--THE END--
