#!/usr/bin/env bash

__file=${BASH_SOURCE[0]}
__dir=$(cd $(dirname ${__file}) && pwd)
__root=$(cd ${__dir}/../ && pwd)
if [[ -f ${__root}/.env ]]; then source ${__root}/.env; fi

run_build_task=${1:-1}
group_name=${2:-${GROUP_NAME}}
webapp_name=${3:-${WEBAPP_NAME}}
acr_name=${4:-${ACR_NAME}}
pg_hostname=${5:-${PG_SERVER_NAME}}
location=${6:-${DEFAULT_LOCATION}}
gh_token=${7:-${GH_TOKEN}}

plan_name="${webapp_name}-plan"

project_name=joshgav/vanpool-manager
image_tag=latest
image_uri=${acr_name}.azurecr.io/${project_name}:${image_tag}

# create group
group_id=$(az group show --name ${group_name} --output tsv --query id 2> /dev/null)
if [[ -z $group_id ]]; then
  group_id=$(az group create --name ${group_name} --location ${location} \
    --output tsv --query id)
fi
echo "group_id: $group_id"

# create registry
acr_id=$(az acr show --name ${acr_name} \
    --resource-group ${group_name} \
    --query 'id' --output tsv 2> /dev/null)
if [[ -z "$acr_id" ]]; then
  acr_id=$(az acr create \
    --name ${acr_name} \
    --resource-group ${group_name} \
    --sku  "Standard" \
    --location ${location} \
    --admin-enabled \
    --output tsv --query 'id')
fi
acr_url=$(az acr show --name ${acr_name} \
    --resource-group ${group_name} \
    --query 'loginServer' --output tsv 2> /dev/null)
echo "acr_id: ${acr_id}"
echo "acr_url: ${acr_url}"

# create container build-task
buildtask_name=buildoncommit
buildtask_id=$(az acr build-task show \
    --name ${buildtask_name} \
    --registry ${acr_name} \
    --resource-group ${group_name} 2> /dev/null)
if [[ -z "$buildtask_id" ]]; then
    buildtask_id=$(az acr build-task create \
        --context "https://github.com/${project_name}" \
        --git-access-token $gh_token \
        --image ${image_uri} \
        --name ${buildtask_name} \
        --registry ${acr_name} \
        --resource-group ${group_name} \
        --commit-trigger-enabled true \
        --output tsv --query id)
fi
if [[ "$run_build_task" == "1" ]]; then
    az acr build-task run --no-logs \
        --name ${buildtask_name} \
        --registry ${acr_name} \
        --resource-group ${group_name}
fi

# create database
pg_rule_name=allow-all
pg_server_id=$(az postgres server show \
    --name ${pg_hostname} --resource-group ${group_name} \
    --output tsv --query id 2> /dev/null)
if [[ -z "$pg_server_id" ]]; then
  # SKUs: https://docs.microsoft.com/en-us/azure/postgresql/concepts-pricing-tiers
  pg_server_id=$(az postgres server create \
    --name ${pg_hostname} \
    --resource-group ${group_name} \
    --location ${location} \
    --admin-user ${POSTGRES_USER} \
    --admin-password ${POSTGRES_PASSWORD} \
    --sku-name 'B_Gen5_2' \
    --ssl-enforcement Disabled \
    --storage-size 51200)

  az postgres server firewall-rule create \
    --name ${pg_rule_name} \
    --resource-group ${group_name} \
    --server-name ${pg_hostname} \
    --start-ip-address '0.0.0.0' \
    --end-ip-address '255.255.255.255' > /dev/null
fi
echo "pg_server_id: ${pg_server_id}"

# create webapp
webapp_id=$(az webapp show --name ${webapp_name} --resource-group ${group_name} \
    --output tsv --query id 2> /dev/null)
if [[ -z "$webapp_id" ]]; then
   plan_id=$(az appservice plan create \
     --name ${plan_name} \
     --resource-group ${group_name} \
     --location ${location} \
     --is-linux \
     --output tsv --query id)
 
   webapp_id=$(az webapp create \
     --name ${webapp_name} \
     --plan ${plan_name} \
     --resource-group ${group_name} \
     --deployment-container-image-name 'scratch' \
     --output tsv --query id)
fi

# configure webapp
az webapp config container set \
 --ids $webapp_id \
 --docker-registry-server-url "https://${acr_url}" \
 --docker-custom-image-name "${acr_url}/${project_name}:${image_tag}"

az webapp config appsettings set \
 --ids $webapp_id \
 --settings \
     "AZURE_CLIENT_ID=${AZURE_CLIENT_ID}" \
     "AZURE_CLIENT_SECRET=${AZURE_CLIENT_SECRET}" \
     "COOKIE_KEY=${COOKIE_KEY}" \
     "POSTGRES_HOSTNAME=${PG_SERVER_NAME}.postgres.database.azure.com" \
     "POSTGRES_PORT=${POSTGRES_PORT}" \
     "POSTGRES_SSLMODE=require" \
     "POSTGRES_DB=postgres" \
     "POSTGRES_USER=${POSTGRES_USER}%40${PG_SERVER_NAME}" \
     "POSTGRES_PASSWORD=${POSTGRES_PASSWORD}" \
     "REDIRECT_HOSTNAME=${webapp_name}.azurewebsites.net"

# create cache
redis_id=$(az redis show --name ${redis_name} --resource-group ${group_name} \
    --output tsv --query id)
if [[ -z $redis_id ]]; then
    redis_id=$(az redis create \
        --name ${redis_name} \
        --resource-group ${group_name} \
        --location ${location} \
        --sku 'Standard' \
        --vm-size 'C4' \
        --output tsv --query id)
fi
echo "redis_id: $redis_id"

# create service bus
# az servicebus create ...

# create event hub
# az eventhubs create

# create functionapp
# az functionapp create ...
