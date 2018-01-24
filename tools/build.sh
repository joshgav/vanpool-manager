#!/usr/bin/env bash
cd ./web
npm install
./node_modules/.bin/webpack

docker run -d --rm \
  --publish 5432:5432 \
  --env-file .env \
  --name postgres_test
  postgres:latest

