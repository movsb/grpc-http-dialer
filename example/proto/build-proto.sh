#!/bin/bash

set -eu

protoc --go_out=plugins=grpc:. ping.proto
