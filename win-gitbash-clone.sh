#! /bin/bash

export DEVEL=/c/tmp/devel
export GOPATH=$DEVEL/negentropia/webserv

[ -d $DEVEL ] || mkdir -p $DEVEL
cd $DEVEL

git clone https://github.com/udhos/negentropia

go_get () {
	local i=$*
	echo go get $i
	go get $i
}

#
# fetch from github with git bash
#
#go_get github.com/vmihailenco/redis
#go_get github.com/vmihailenco/redis/v2
go_get gopkg.in/redis.v2 ;# github: github.com/go-redis/redis
#go_get github.com/HairyMezican/goauth2/oauth
#go_get github.com/spate/vectormath
go_get github.com/udhos/vectormath

#
# fetch from google code (could be under cmd prompt)
#
#go_get code.google.com/p/go.net/websocket
go_get golang.org/x/net/websocket
#go_get code.google.com/p/goauth2/oauth
go_get github.com/udhos/goauth2/oauth

# gopherjs dependencies
cat $DEVEL/negentropia/gopherjs-deps.txt | while read i; do
	j=`echo $i` ;# withstand CR LR issues
	go_get $j
done
