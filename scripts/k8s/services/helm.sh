#!/usr/bin/env bash

this_dir=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
k8s_dir=$(dirname "${this_dir}")
root_dir=$(dirname "${k8s_dir}/..")
source ${k8s_dir}/common.sh
source ${k8s_dir}/vars.sh

if [[ ! $(which jq) ]]; then >&2 echo "first install jq"; exit 2; fi

install-helm () {
	curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get > get_helm.sh
	chmod +x get_helm.sh
	./get_helm.sh
  rm ./get_helm.sh
}

install-tiller () {
	kubectl apply -f "${this_dir}/helm-rbac.yaml"
	helm init --upgrade --service-account tiller
}

install-helm
install-tiller
