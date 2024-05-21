#!/bin/sh

CURDIR=$(cd $(dirname $0); pwd)

exec "$CURDIR/bin/tokengateway.api" $*
