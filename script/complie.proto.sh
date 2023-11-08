#!/bin/bash
cd internal/model
protoc -I=. --go_out=. ./*.proto
cd ../..

