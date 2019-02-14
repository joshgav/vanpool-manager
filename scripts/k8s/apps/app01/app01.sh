# build and deploy a simple Node.js app to K8s
#!/usr/bin/env bash
scripts_dir=$(dirname $(cd $(dirname "${BASH_SOURCE[0]}") && pwd))
root_dir=$(dirname "${scripts_dir}")
source ${scripts_dir}/common.sh
source ${scripts_dir}/vars.sh

app_name=app01
task_name=build-${app_name}
task_id=$(az acr task show \
    --name ${task_name} \
    --registry ${registry_name} \
    --output tsv --query id)
if [[ -z "$task_id" ]]; then
    if [[ -z "$gh_token" ]]; then echo "set \$gh_token first"; fi
    task_id=$(az acr task create \
        --registry ${registry_name} \
        --name ${task_name} \
        --context "https://github.com/joshgav/node-scratch.git" \
        --file "Dockerfile" \
        --git-access-token ${gh_token} \
        --image 'node-scratch:latest' \
        --output tsv --query id)
fi
info "ensured task_id: ${task_id}"

# TODO: check if a run is necessary
# az acr task run \
#     --name ${task_name} \
#     --registry ${registry_name}

kubectl apply -f "${scripts_dir}/apps/${app_name}.yaml"
external_ip=$(kubectl get services "${app_name}-service" -o json |
    jq --raw-output '.status.loadBalancer.ingress[0].ip')
info "access at http://${external_ip}"