#! /bin/bash

export DEVEL=/c/tmp/devel
export GOPATH=$DEVEL/negentropia/webserv

[ -d $DEVEL ] || mkdir -p $DEVEL
cd $DEVEL

git clone https://code.google.com/p/negentropia/

go_get () {
	local i=$1
	echo go get $i
	go get $i
}

# fetch from github with git bash
#go_get github.com/vmihailenco/redis
go_get github.com/vmihailenco/redis/v2
go_get github.com/HairyMezican/goauth2/oauth
#go_get github.com/spate/vectormath
go_get github.com/udhos/vectormath

# fetch from google code (could be under cmd prompt)
go_get code.google.com/p/go.net/websocket
