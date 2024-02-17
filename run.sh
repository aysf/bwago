#!/bin/bash

source env.sh

go build -o bwago cmd/web/*.go && ./bwago 