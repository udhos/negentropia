
REQUIREMENTS

# Install Go Language
# http://golang.org/doc/install

# Install GIT - required to clone the source tree
# http://git-scm.com/downloads

# Install Mercurial - required for packages such as: go get code.google.com/p/goauth2/oauth
# http://mercurial.selenic.com/downloads/

# Install Redis
# http://redis.io/download
#
# Microsoft Official Redis for Windows at Microsoft Open Tech:
# https://github.com/MSOpenTech/redis
# https://github.com/MSOpenTech/redis/blob/bksavecow/msvs/bin/release/redisbin.zip
#
# Redis 2.4.5 for Windows 7 64-bit:
# https://github.com/dmajkic/redis/downloads


GENERAL BUILDING INSTRUCTIONS:

# 1. Clone the git repository:
git clone https://code.google.com/p/negentropia/

# 2. Set GOPATH to negentropia\webserv

# 3. install goauth2
go get code.google.com/p/goauth2/oauth

# 4. Install redis client library
go get github.com/vmihailenco/redis

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

go get github.com/vmihailenco/redis

## windows dos prompt:

@rem start redis
@rem dmajkic redis:
c:\redis-2.4.5-win32-win64\64bit\redis-server.exe
@rem microsoft redis:
c:\redisbin\redis-server.exe

set DEVEL=C:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem install goauth2
go get code.google.com/p/goauth2/oauth

@rem build
go install negentropia\webserv

@rem run:
@rem   -- google login requires Google API "Client ID" and "Client secret"
%DEVEL%\negentropia\webserv\bin\webserv.exe -gId=putIdHere -gSecret=putSecretHere
@rem   -- if you don't need google login:
%DEVEL%\negentropia\webserv\bin\webserv.exe


RUNNING / TESTING

Open http://localhost:8080/n/

--THE END--
