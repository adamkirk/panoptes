#!/bin/bash

# If a command was passed to the container run the binary 
if [ "$#" != "0" ]; then
    go run ./cmd/main.go "$@"
    exit "$?"
fi

DEBUG_OPT=""

if [ "$AIR_DEBUG" == "true" ]; then
    DEBUG_OPT="-d"
fi

air $DEBUG_OPT -build.bin ./build/heimdallr -build.cmd "./hack/build.sh"