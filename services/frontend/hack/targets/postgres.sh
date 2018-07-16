#!/usr/bin/env bash
set -o errexit -o pipefail
__DIRNAME=$(dirname $(realpath "${BASH_SOURCE[0]}"))
__ROOTDIR=$(dirname "$(dirname "$__DIRNAME")")

# $1: command: ensure | restart | stop | connect
# $2: env: local | azure
function postgres {
    source ${__ROOTDIR}/.env
    declare command=${1:-"ensure"}
    declare env=${2:-"local"}
    declare group_name=$(printf '%s-g' "${AZURE_APP}")
    declare vnet_name=$(printf '%s-vnet' "${AZURE_APP}")
    declare main_subnet_name=$(printf '%s-main_subnet' "${AZURE_APP}")

    DB_IMAGE_TAG=library/postgres:alpine
    DB_CONTAINER_NAME=buffalo-postgres-01

    # local
    if [[ "$env" == "local" ]]; then
        echo "setup postgres locally with official Docker image"
        docker image pull $DB_IMAGE_TAG

        # check if we already have a running container
        CID=$(docker container ls \
            --filter "name=${DB_CONTAINER_NAME}$" --quiet)
        RUNNING=$(if [[ -n "$CID" ]]; then echo 1; else echo 0; fi)

        if [[ -n "$CID" ]]; then
            echo "Container [${DB_CONTAINER_NAME}] already running [CID: ${CID}]"
            docker container ls \
                --filter "id=$CID" \
                --format "Status: {{ .Status }}"
            if [[ ("$command" == "restart") || ("$command" == "stop") ]]; then
                docker kill $CID # && docker rm $CID 
            fi
        fi


        if [[ ("$command" == "ensure" && -z "$CID") || \
              ("$command" == "connect" && -z "$CID") || \
              ("$command" == "restart") \
           ]]; then

            docker container run \
                --name $DB_CONTAINER_NAME \
                --detach \
                --rm \
                --publish 5432:5432 \
                --env-file ${__ROOTDIR}/.env \
                --mount $(printf "%s,%s,%s" \
                    "type=bind" \
                    "source=${__DIRNAME}/../testdata/pgdata" \
                    "target=/var/lib/postgresql/data") \
                $DB_IMAGE_TAG
        fi

        if [[ ("$command" == "connect") ]]; then
            docker run \
                --rm -it \
                --link ${DB_CONTAINER_NAME}:postgres \
                ${DB_IMAGE_TAG} \
                psql -h postgres \
                -d ${POSTGRES_DB} \
                -U ${POSTGRES_USER}
        fi # //end "connect"
    fi # // "$env" == "local"

    # azure

    if [[ "$env" == "azure" ]]; then

        echo "azure setup"

        set +e
        PG_SERVER_ID=$(az postgres server show \
            --name ${AZURE_POSTGRES_HOST} \
            --resource-group $group_name \
            --output tsv --query id 2> /dev/null)
        if [[ -n "$PG_SERVER_ID" ]]; then echo found server: [$PG_SERVER_ID]; fi
        set -e

        if [[ "$command" == "stop" || "$command" == "restart" ]]; then
            echo "stop | restart"
            if [[ -n "$PG_SERVER_ID" ]]; then
                az postgres server delete --yes --ids $PG_SERVER_ID 
                # blank this so it will be recreated later
                PG_SERVER_ID=
            fi
        fi

        if [[ "$command" == "ensure" || "$command" == "restart" ||
              "$command" == "connect" ]]; then

            echo "ensure|restart|connect,azure"
            echo "ensuring group, vnet and subnet"
            ${__DIRNAME}/vnet.sh

            # check again cause postgres group could be different than vnet group
            PG_GROUP_ID=$(az group show \
                --name $group_name \
                --output tsv --query id)
            if [[ -n "$PG_GROUP_ID" ]]; then echo found group: [$PG_GROUP_ID]; fi

            if [[ -z "$PG_GROUP_ID" ]]; then
                PG_GROUP_ID=$(az group create \
                    --name $group_name \
                    --location ${AZURE_POSTGRES_LOCATION})
                echo created group: $PG_GROUP_ID
            fi

            if [[ -z "$PG_SERVER_ID" ]]; then
                PG_SERVER_ID=$(az postgres server create \
                    --name ${AZURE_POSTGRES_HOST} \
                    --resource-group $group_name \
                    --sku-name ${AZURE_POSTGRES_SKU} \
                    --location ${AZURE_POSTGRES_LOCATION} \
                    --ssl-enforcement "Enabled" \
                    --storage-size 10240 \
                    --version ${AZURE_POSTGRES_VERSION} \
                    --admin-user ${POSTGRES_USER} \
                    --admin-password ${POSTGRES_PASSWORD} \
                    --output tsv --query id)
                echo created pg_server: [$PG_SERVER_ID]
                

                RULE_ID=$(az postgres server vnet-rule create \
                    --name 'allow-pg-to-main_subnet' \
                    --resource-group $group_name \
                    --server-name ${AZURE_POSTGRES_HOST} \
                    --subnet $main_subnet_name \
                    --vnet-name $vnet_name \
                    --output tsv --query id)
                echo created vnet-rule: [$RULE_ID]

                DB_ID=$(az postgres db create \
                    --name ${POSTGRES_DB} \
                    --resource-group $group_name \
                    --server-name ${AZURE_POSTGRES_HOST} \
                    --output tsv --query id)
                echo created db: [$DB_ID]
            fi
        fi # //end "ensure|restart|connect"
    fi # // "$env" == "azure"

    if [[ "$command" == "connect" ]]; then
        azure_host=$(printf \
            '%s.%s' ${AZURE_POSTGRES_HOST} "postgres.database.azure.com")
        azure_user=$(printf \
            '%s@%s' ${POSTGRES_USER} ${AZURE_POSTGRES_HOST})
        my_ip_address=$(curl -sS ifconfig.co)

        echo "enabling temporary access from $my_ip_address"
        FWRULE_ID=$(az postgres server firewall-rule create \
            --name 'temp-psql-access' \
            --resource-group $group_name \
            --server-name ${AZURE_POSTGRES_HOST} \
            --start-ip-address $my_ip_address \
            --end-ip-address $my_ip_address \
            --output tsv --query id)

        echo "psql -h $azure_host -U $azure_user -d ${POSTGRES_DB} -v \"sslmode=true\""
        psql \
            -h "$azure_host" \
            -U "$azure_user" \
            -d ${POSTGRES_DB} \
            -v "sslmode=true"

        az postgres server firewall-rule delete --ids $FWRULE_ID --yes
    fi
}

postgres "$@"
