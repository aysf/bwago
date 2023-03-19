#!/bin/bash

source env.sh

go build -o bwago cmd/web/*.go && ./bwago \
    --cache=$BWAGO_CACHE \
    --dbname=$BWAGO_DB_NAME \
    --dbuser=$BWAGO_DB_USER \
    --dbpass=$BWAGO_DB_PASS