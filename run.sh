#!/bin/bash
trap "rm main;kill 0" EXIT

go build -o main
./main -port=8001 &
./main -port=8002 &
./main -port=8003 -api=1 &
sleep 2
echo ">>>> start test"
curl "http://127.0.0.1:9999/api?key=Tom" &
curl "http://127.0.0.1:9999/api?key=Tom"

wait

