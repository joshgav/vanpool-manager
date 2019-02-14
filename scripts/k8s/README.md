# azure-dapp

Deployment and configuration artifacts for a distributed app on Azure.  This app uses:

* Azure Kubernetes Service (AKS) to schedule workloads.
* Azure Container Registry (ACR) to store artifacts (images).
* Helm to manage workload lifecycle.
* Service Catalog and Azure Service Broker to provision managed services.
* Azure AD for identifying and authenticating users.
* Azure PostgreSQL Database for storing persistent data.
* Azure App Insights to collect traces, logs and metrics.

## Run

1. Set insecure variables in `./scripts/k8s/vars.sh`. Set secure variables in a
   copy of `.env.tpl`.
1. Create AKS and ACR resources with `./scripts/k8s/cluster.sh`.
1. Install Helm and Tiller with `./scripts/k8s/services.helm.sh`
1. Install Service Catalog and Azure Service Broker with
   `./scripts/services/svccat.sh`.
1. Create a principal and secret for OIDC with `./scripts/services/oidc.sh`.
    * The principal must be created in the portal if MSA auth is also desired.
1. Set up ACR credentials for use in K8s with `./scripts/acr_creds.sh`.
1. Create the vanpool-manager app with `./scripts/vanpool-manager.sh`.
    * `gh_token` must be set to connect the GitHub repo to ACR.

# License

See [LICENSE.md](./LICENSE.md).
