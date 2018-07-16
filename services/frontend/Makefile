azure-infra:
	./hack/targets/vnet.sh ensure
	./hack/targets/postgres.sh ensure azure

local-infra:
	./hack/targets/postgres.sh ensure local

.PHONY: azure-infra local-infra
