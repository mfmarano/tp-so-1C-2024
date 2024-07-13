#!/bin/bash

if [ -z "$KERNEL_PORT" ]; then
    echo "The KERNEL_PORT is not set"
    echo "Using default port 8080"
    KERNEL_PORT=8001
fi

if [ -z "$KERNEL_HOST" ]; then
    echo "The KERNEL_HOST is not set"
    echo "Using default host localhost"
    KERNEL_HOST=127.0.0.1
fi

curl --location --request PUT "http://$KERNEL_HOST:$KERNEL_PORT/process" \
--header 'Content-Type: application/json' \
--data '{
    "pid": 1,
    "path": "/scripts_memoria/FS_1"
}'

curl --location --request PUT "http://$KERNEL_HOST:$KERNEL_PORT/process" \
--header 'Content-Type: application/json' \
--data '{
    "pid": 2,
    "path": "/scripts_memoria/FS_2"
}'