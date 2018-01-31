#!/usr/bin/env bash
set -o allexport
source .env
set +o allexport

export AZ_LOCATION=westus
export AZ_GROUP_NAME=vanpool-manager-dev
export AZ_ACR_NAME=vanpoolmanagerdev
export AZ_POSTGRES_HOSTNAME=vanpooldb
export AZ_POSTGRES_FWNAME=pgfw01
export AZ_WEBAPP_NAME=vanpool-manager-dev
export AZ_WEBAPP_PLAN_NAME=vpmgrplan01

export tag=${AZ_ACR_NAME}.azurecr.io/joshgav/vanpool-manager:latest

# create resource group
export group_exists=$(az group exists --name ${AZ_GROUP_NAME} --output json)
if [ $group_exists = "false" ]; then
  az group create --name ${AZ_GROUP_NAME} --location ${AZ_LOCATION}
fi

# create ACR registry
az acr show --name ${AZ_ACR_NAME} 2>&1 > /dev/null
export acr_not_exists=$(echo $?)
if [ "x$acr_not_exists" = "x1" ]; then
  az acr create --verbose \
    --name ${AZ_ACR_NAME} \
    --resource-group ${AZ_GROUP_NAME} \
    --sku  "Standard" \
    --location ${AZ_LOCATION} \
    --admin-enabled
fi

# push image to registry
docker build -t $tag .
az acr login --name ${AZ_ACR_NAME} --resource-group ${AZ_GROUP_NAME}
docker push $tag

# postgres DB
az postgres server show --name ${AZ_POSTGRES_HOSTNAME} --resource-group ${AZ_GROUP_NAME} 2> /dev/null
export pg_not_exists=$(echo $?)
if [ "x$pg_not_exists" = "x1" ]; then
  az postgres server create \
    --name ${AZ_POSTGRES_HOSTNAME} \
    --resource-group ${AZ_GROUP_NAME} \
    --location ${AZ_LOCATION} \
    --admin-user ${POSTGRES_USER} \
    --admin-password ${POSTGRES_PASSWORD} \
    --performance-tier Basic --compute-units 100 --ssl-enforcement Disabled --storage-size 51200

  az postgres server firewall-rule create \
    --name ${AZ_POSTGRES_FWNAME} \
    --resource-group ${AZ_GROUP_NAME} \
    --server-name ${AZ_POSTGRES_HOSTNAME} \
    --start-ip-address '0.0.0.0' \
    --end-ip-address '255.255.255.255'
fi

# bug with empty return at the moment
# az webapp show --name ${AZ_WEBAPP_NAME} --resource-group ${AZ_GROUP_NAME} 2> /dev/null
# export webapp_not_exists=$(echo $?)

# export webapp_not_exists=1
# if [ "x$webapp_not_exists" = "x1" ]; then
#   az appservice plan create \
#     --name ${AZ_WEBAPP_PLAN_NAME} \
#     --resource-group ${AZ_GROUP_NAME} \
#     --location ${AZ_LOCATION} \
#     --is-linux
# 
#   az webapp create \
#     --name ${AZ_WEBAPP_NAME} \
#     --plan ${AZ_WEBAPP_PLAN_NAME} \
#     --resource-group ${AZ_GROUP_NAME} 
# TODO: fix deployment of container
# fi

