#!/usr/bin/env bash

if [ ! -f install.sh ]; then
echo 'install must be run within its container folder' 1>&2
exit 1
fi
cd ..
CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"
export GOBIN="$CURDIR"/bin
echo "current gopath is "  $GOPATH
gofmt -w src

go install server

export GOPATH="$OLDGOPATH"

cd bin
mv server go-sms-server


echo 'finished'