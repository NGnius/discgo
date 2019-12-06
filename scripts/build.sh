#!/bin/bash
# Assumption: script is run from project root
cd ./discgo
go mod vendor
go build
cd ..
