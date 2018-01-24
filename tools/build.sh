#!/usr/bin/env bash
cd ./web
npm install
npm run build
cd ..

docker run -d --rm \
  --publish 5432:5432 \
  --env-file .env \
  --name postgres_test \
  postgres:latest

sleep 3 # give the db a chance to start twice

go test -v github.com/joshgav/go-demo/data

docker container stop postgres_test
