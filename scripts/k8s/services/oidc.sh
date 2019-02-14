#!/usr/bin/env bash
this_dir=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
scripts_dir=$(dirname ${this_dir})
source ${scripts_dir}/common.sh
source ${scripts_dir}/vars.sh

# az CLI doesn't yet support creation of principals
# with ability to authenticate MSA accounts.
# Set `oidc_client_id` and `oidc_client_secret` via vars.sh
# to bypass automatic principal creation
client_id=${oidc_client_id}
client_secret=${oidc_client_secret}

# app URLs
app_name=vanpool-manager
app_domain=microsoft.com
app_scheme=https
app_hostname="${app_name}.${app_domain}"
app_identifier="${app_scheme}://${app_hostname}"

app_url="${app_scheme}://${app_name}.${app_domain}"
declare -a app_reply_urls=(
  "http://localhost:8080/login"
  "${app_url}/login"
)
# /end app URLs

if [[ -z "${client_id}" ]]; then
  info "OIDC client ID not provided; creating new"

  az ad app create \
    --display-name ${app_name} \
    --identifier-uris ${app_identifier} \
    --reply-urls "${app_reply_urls[@]}" \
    --available-to-other-tenants true \
    --key-type password \
    --required-resource-accesses @"${this_dir}/oidc_manifest.json" \
    1> /dev/null

  # AzureADMyOrg, AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount
  # doesn't work
  # az ad app update \
  #   --id ${app_identifier} \
  #   --set '.signInAudience=AzureADandPersonalMicrosoftAccount'

  az ad sp create --id $app_identifier 1> /dev/null
  json=$(az ad app credential reset --id $app_identifier --output json)

  client_id=$(echo $json | jq --raw-output ".appId")
  client_secret=$(echo $json | jq --raw-output ".password")
else
  info "using provided OIDC client ID [${client_id}]"
fi

identity_name=${app_name}-identity
kubectl delete secret ${identity_name}
kubectl create secret generic ${identity_name} \
  --from-literal=client_id=${client_id} \
  --from-literal=client_secret=${client_secret} \
  --from-literal=redirect_hostname=${app_hostname} \
  --from-literal=redirect_scheme=${app_scheme}
