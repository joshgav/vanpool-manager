#!/usr/bin/env bash
scripts_dir=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
source ${scripts_dir}/common.sh
source ${scripts_dir}/vars.sh

set-registry-creds () {
  local registry_name=${1:-${registry_name}}

  az acr update --name $registry_name \
    --admin-enabled true

  json=$(az acr credential show --name $registry_name --output json)
  docker_server=$(az acr show \
      --name $registry_name --query "loginServer" --output tsv)
  docker_username=$(echo $json | jq --raw-output ".username")
  docker_password=$(echo $json | jq --raw-output ".passwords[0].value")
  docker_email="nobody@microsoft.com"

  kubectl create secret docker-registry acr-cred \
    --docker-server=${docker_server} \
    --docker-username=${docker_username} \
    --docker-password=${docker_password} \
    --docker-email=${docker_email}
}

set-registry-creds
