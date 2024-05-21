#!/usr/bin/env bash

RUN_NAME="tokengateway.api"
mkdir -p output/bin output/conf
cp script/* output/
cp conf/* output/conf/
chmod +x output/bootstrap.sh

TAG='musl'
if [[ `uname` == 'Darwin' ]]; then
	TAG='dynamic'
fi

go build -tags=$TAG,jsoniter -o output/bin/${RUN_NAME}