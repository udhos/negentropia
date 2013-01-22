
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
#
# http://code.google.com/p/goauth2/					OAuth 2.0 for Go clients. Doesn't work with Facebook: http://code.google.com/p/goauth2/issues/detail?id=4 
# https://github.com/robfig/goauth2					A fork of code.google.com/p/goauth2 that supports Facebook
# https://github.com/HairyMezican/goauth2/			This is mostly copied from http://code.google.com/p/goauth2/
#													The original code will fail when contacting facebook; this code fixes that problem
# http://code.google.com/r/jasonmcvetta-goauth2/	This clone contains changes to stock Goauth2 detailed by Ryan.C.K. here: http://code.google.com/p/goauth2/issues/detail?id=4

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

# fetch from github with git bash
go get github.com/vmihailenco/redis
go get github.com/HairyMezican/goauth2/oauth

## windows dos prompt:

@rem start redis
@rem dmajkic redis:
c:\redis-2.4.5-win32-win64\64bit\redis-server.exe
@rem microsoft redis:
c:\redisbin\redis-server.exe

set DEVEL=C:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

@rem fetch from code.google.com with DOS prompt
@rem
@rem install goauth2
@rem
@rem facebook broken:
@rem go get code.google.com/p/goauth2/oauth
@rem
@rem google broken:
@rem go get github.com/robfig/goauth2/oauth
@rem
@rem go get broken:
@rem go get code.google.com/r/jasonmcvetta-goauth2/
@rem
@rem load from git bash:
@rem go get github.com/HairyMezican/goauth2/oauth

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
