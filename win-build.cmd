@rem win-build

set DEVEL=C:\tmp\devel
set GOPATH=%DEVEL%\negentropia\webserv

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
go install negentropia/webserv
go install negentropia/world

@rem eof
