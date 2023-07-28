#!/usr/bin/env bash

protoc -I=./ --go_out=./ ctrl.proto client.proto
