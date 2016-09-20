#!/bin/bash

current=`pwd`
go get github.com/landru29/swaggo
cd $GOPATH/src/github.com/landru29/swaggo
go build
go install
cd $current
