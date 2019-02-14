#!/usr/bin/env bash

this_dir=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
k8s_dir=$(dirname "${this_dir}")
root_dir=$(dirname "${k8s_dir}/..")
source ${k8s_dir}/common.sh
source ${k8s_dir}/vars.sh

if [[ ! $(which jq) ]]; then >&2 echo "first install jq"; exit 2; fi

install-azurebroker () {
	local sp_name=${1:-"service-catalog-azure-broker"}

	helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
  helm install svc-cat/catalog \
	  --name catalog --namespace catalog \
    --set apiserver.storage.etcd.persistence.enabled=true

  kubectl apply -f ${this_dir}/svccat-role-fix.yaml

	sp_name="http://${sp_name}"
	sub_id=$(az account show --query id --output tsv)
  # TODO: use existing or reset password if principal exists
	json=$(az ad sp create-for-rbac \
		--name ${sp_name} \
		--scopes "/subscriptions/${sub_id}" \
		--output json)
	client_id=$(echo $json | jq .appId)
	password=$(echo $json | jq .password)
	tenant_id=$(echo $json | jq .tenant)

	helm repo add azure https://kubernetescharts.blob.core.windows.net/azure
	helm install azure/open-service-broker-azure \
		--set modules.minStability=experimental \
		--name azure-broker --namespace catalog \
		--set azure.clientId=${client_id} \
		--set azure.clientSecret=${password} \
		--set azure.subscriptionId=${sub_id} \
		--set azure.tenantId=${tenant_id}
}

install-azurebroker
