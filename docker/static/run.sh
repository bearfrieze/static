#!/bin/bash
NAME=static-dev-static
NIMBUS=nimbus-dev-nimbus
WORKDIR=/go/src/github.com/bearfrieze/static
PORT=8081
docker run --rm --name $NAME -v $HOME/go:/go -w $WORKDIR -p $PORT:$PORT -it --link $NIMBUS:nimbus $NAME