#!/usr/bin/env bash
scripts_dir=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
root_dir=$(dirname "${scripts_dir}")
source ${scripts_dir}/common.sh
source ${scripts_dir}/vars.sh

# this script deploys the following:
#   group
#   service principal
#   AKS cluster
#   ACR registry
#   App Insights sink

## ensure group
group_id=$(ensure_group $group_name $group_location)
info "ensured group_id: $group_id"
# /end ensure group

## ensure service principal
app_id=$(az ad sp show --id "${cluster_identity}" \
    --output tsv --query appId 2> /dev/null)
if [[ -z "$app_id" ]]; then
    app_id=$(az ad sp create-for-rbac \
        --name "${cluster_identity}" \
        --skip-assignment true \
        --scopes ${group_id} \
        --output tsv --query appId)
fi
info "ensured app_id: ${app_id}"
# /end ensure service principal

## ensure AKS cluster
aks_id=$(az aks show \
    --name ${cluster_name} \
    --resource-group ${cluster_group_name} \
    --output tsv --query id)
if [[ -z "$aks_id" ]]; then
    info "cluster not found; starting cluster create"
    info "first resetting sp credential to reuse"
    app_password=$(az ad sp credential reset \
        --name "${cluster_identity}" \
        --output tsv --query password)

    info 'calling `az aks create`'
    az aks create \
        --name ${cluster_name} \
        --resource-group ${cluster_group_name} \
        --location ${cluster_location} \
        --dns-name-prefix ${cluster_prefix} \
        --service-principal ${app_id} \
        --client-secret ${app_password} \
        --network-plugin kubenet \
        --ssh-key-value ${ssh_pubkey_path} \
        --no-wait

    info "awaiting cluster creation"
    az aks wait \
        --name ${cluster_name} \
        --resource-group ${cluster_group_name} \
        --created
fi

info "setting AKS credentials in ~/.kube/config"
az aks get-credentials \
    --name ${cluster_name} \
    --resource-group ${cluster_group_name} \
    --admin >> /dev/null

info "ensured aks_id: ${aks_id}"
# /end ensure AKS cluster

## ensure container registry
acr_id=$(az acr show \
    --name ${registry_name} \
    --resource-group ${registry_group} \
    --output tsv --query id)
if [[ -z "$acr_id" ]]; then
    info "creating container registry"
    acr_id=$(az acr create \
        --name ${registry_name} \
        --resource-group ${registry_group} \
        --location ${registry_location} \
        --sku Standard \
        --output tsv --query id)
    
    # grant AKS principal read rights to ACR
    role_assignment_id=$(az role assignment create \
        --assignee ${app_id} \
        --scope ${acr_id} \
        --role Reader \
        --output tsv --query id)
    
    # add ACR credentials to ~/.docker
    az acr login \
        --name ${registry_name} \
        --resource-group ${registry_group} >> /dev/null
fi
info "ensured acr_id: ${acr_id}"
# /end ensure container registry

## ensure app insights
appinsights_id=$(az resource show \
    --name ${appinsights_resource_name} \
    --resource-group ${appinsights_group} \
    --resource-type "Microsoft.Insights/components" \
    --query id --output tsv)
if [[ -z "$appinsights_id" ]]; then
    appinsights_id=$(az resource create \
        --resource-group ${appinsights_group} \
        --resource-type "Microsoft.Insights/components" \
        --name ${appinsights_resource_name} \
        --location ${appinsights_location} \
        --properties "{ \"Application_Type\":\"other\", \"Flow_Type\":\"Redfield\" }" \
        --output tsv --query id)
fi
appinsights_ikey=$(az resource show \
    --name ${appinsights_resource_name} \
    --resource-group ${appinsights_group} \
    --resource-type "Microsoft.Insights/components" \
    --query properties.InstrumentationKey --output tsv)
info "ensured appinsights_id: ${appinsights_id} with ikey: ${appinsights_ikey}"
# /end ensure app insights
