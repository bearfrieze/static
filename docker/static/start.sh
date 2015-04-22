#!/bin/bash
export NIMBUS=http://$NIMBUS_PORT_8080_TCP_ADDR:$NIMBUS_PORT_8080_TCP_PORT
go run main.go