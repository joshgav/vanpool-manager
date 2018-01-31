#!/usr/bin/env bash
export IMAGE_TAG=joshgav/vanpool-manager:latest
export PG_CONTAINER_NAME=test_postgres_db
export API_CONTAINER_NAME=test_riders_api
export PACKAGE_NAME=joshgav/go-demo

docker container stop ${API_CONTAINER_NAME}
docker container stop ${PG_CONTAINER_NAME}

docker build --tag ${IMAGE_TAG} .

docker run --detach --rm \
  --publish 5432:5432 \
  --env-file .env \
  --name ${PG_CONTAINER_NAME} \
  postgres:latest

# give db chance to initialize
sleep 3

docker run --detach --rm \
  --publish 8080:8080 \
  --env-file .env \
  --name ${API_CONTAINER_NAME} \
  ${IMAGE_TAG}

