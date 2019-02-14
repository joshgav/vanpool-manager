vars_dir=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
root_dir=$(cd "${vars_dir}/../.." && pwd)
if [[ -f "${root_dir}/.env" ]]; then source "${root_dir}/.env"; fi

# cluster vars
group_name=joshgav-k8s-master-group
group_location=westus2

cluster_name=joshgav-k8s-01
cluster_identity="https://${cluster_name}"
cluster_group_name=${group_name}
cluster_location=${group_location}
cluster_prefix=joshgav-k8s-01
ssh_pubkey_path="~/.ssh/id_rsa.pub"

registry_name=joshgavhub
registry_group=${group_name}
registry_location=${group_location}

appinsights_resource_name=joshgav-k8s-01-insights
appinsights_group=${group_name}
appinsights_location=${group_location}

# used to connect GitHub repo to ACR
# `GH_TOKEN` should be specified in `.env`
gh_token=${GH_TOKEN}

# created in portal, vars set in .env
oidc_client_id=${OIDC_CLIENT_ID}
oidc_client_secret=${OIDC_CLIENT_SECRET}
