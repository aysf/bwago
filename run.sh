#!/bin/bash

source env.sh

go build -o bwago cmd/web/*.go && ./bwago --cache=false --dbname=bookings --dbuser=postgres --dbpass=$DB_PASS