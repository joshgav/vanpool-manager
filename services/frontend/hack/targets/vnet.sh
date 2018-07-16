#!/usr/bin/env bash
set -o errexit -o pipefail
__DIRNAME=$(dirname $(realpath "${BASH_SOURCE[0]}"))
__ROOTDIR=$(dirname "$(dirname "$__DIRNAME")")

function vnet {
    source ${__ROOTDIR}/.env
    local command=${1:-"ensure"}
    local group_name=${2:-"$(printf '%s-g' ${AZURE_APP})"}
    local vnet_name=${3:-"$(printf '%s-vnet' ${AZURE_APP})"}
    local main_subnet_name=${4:-"$(printf '%s-main_subnet' ${AZURE_APP})"}
    local vnet_addresses=${5:-"10.0.0.0/8"}
    local subnet_addresses=${6:-"10.0.0.0/16"}
    local location=${7:-${AZURE_LOCATION_DEFAULT}}
    local tags="vanpool-manager"

    GROUP_ID=$(az group show \
        --name $group_name \
        --output tsv --query id 2> /dev/null)
    if [[ -n "$GROUP_ID" ]]; then echo found group: [$GROUP_ID]; fi

    if [[ -z "$GROUP_ID" ]]; then
        GROUP_ID=$(az group create \
            --name $group_name
            --location $location)
        echo created group: [$GROUP_ID]
    fi

    VNET_ID=$(az network vnet show \
        --name $vnet_name \
        --resource-group $group_name \
        --output tsv --query id)
    if [[ -n "$VNET_ID" ]]; then echo found vnet: [$VNET_ID]; fi

    if [[ -z "$VNET_ID" ]]; then
        VNET_ID=$(az network vnet create \
            --name $vnet_name \
            --resource-group $group_name \
            --address-prefixes '10.0.0.0/8' \
            --location $location \
            --output tsv --query id)
        echo created vnet: [$VNET_ID]
    fi
 
     # not all are supported in all locations
     declare -a vnet_service_endpoints=( \
        'Microsoft.Storage' \
        'Microsoft.Sql' \
        'Microsoft.AzureActiveDirectory' \
        'Microsoft.AzureCosmosDB' \
        'Microsoft.Web' \
        'Microsoft.NetworkServiceEndpointTest' \
        'Microsoft.KeyVault' \
        'Microsoft.EventHub' \
        'Microsoft.ServiceBus' \
        'Microsoft.ContainerRegistry')

    MAIN_SUBNET_ID=$(az network vnet subnet show \
        --name $main_subnet_name \
        --vnet-name $vnet_name \
        --resource-group $group_name \
        --output tsv --query id)
    if [[ -n "$MAIN_SUBNET_ID" ]]; then echo found subnet: [$MAIN_SUBNET_ID]; fi

    if [[ -z "$MAIN_SUBNET_ID" ]]; then
        MAIN_SUBNET_ID=$(az network vnet subnet create \
            --name $main_subnet_name \
            --vnet-name $vnet_name \
            --resource-group $group_name \
            --service-endpoints "Microsoft.Sql" \
            --address-prefix '10.0.0.0/16')
        echo created main_subnet: [$MAIN_SUBNET_ID]
    fi
}

vnet "$@"
