#! /bin/bash

export DEVEL=/c/tmp/devel
export GOPATH=$DEVEL/negentropia/webserv

mkdir -p $DEVEL
cd $DEVEL

git clone https://code.google.com/p/negentropia/

# fetch from github thru git bash
go get github.com/vmihailenco/redis
go get github.com/HairyMezican/goauth2/oauth

# fetch from google code (could be under cmd prompt)
go get code.google.com/p/go.net/websocket