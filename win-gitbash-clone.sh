#! /bin/bash

export DEVEL=/c/tmp/devel
export GOPATH=$DEVEL/negentropia/webserv

mkdir -p $DEVEL
cd $DEVEL

git clone https://code.google.com/p/negentropia/

go get github.com/vmihailenco/redis
go get github.com/HairyMezican/goauth2/oauth