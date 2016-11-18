#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
GRAY='\033[0;37m'
NC='\033[0m' # No Color

# GOPATH=`godep path`:$GOPATH

if [ $1 == "run" ]; then
  go build -o bin/avatar-go ./ && ./bin/avatar-go --debug
  exit
fi

if [ $1 == "config" ]; then
  go build -o bin/avatar-config ./cmd/avatar-config && ./bin/avatar-config
  exit
fi

if [ $1 == "setdefault" ]; then
  go build -o bin/avatar-config ./cmd/avatar-config && ./bin/avatar-config setdefault $2
  exit
fi

if [ $1 == "build" ]; then
  go build -a -o bin/avatar-config ./cmd/avatar-config
  go build -a -o bin/avatar-go ./
  if [ "$?" = "0" ]; then
    echo -e "${GREEN}Completed all builds.${NC}"
  else
    echo -e "${RED}Build failed, please check errors.${NC}"
  fi
  exit
fi

echo "Please include \"config\", \"run\" or \"build\" options."
