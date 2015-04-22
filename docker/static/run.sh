#!/bin/bash
NAME=static-dev-static
WORKDIR=/go/src/github.com/bearfrieze/static
PORT=80
docker run --rm --name $NAME -v $HOME/go:/go -w $WORKDIR -p $PORT:$PORT -it $NAME