# `info` prefixes $1 and writes to stderr
# purpose is to avoid interleaving logs in stdout
function info() {
    local message=$1
    >&2 echo "[info $(date --iso-8601=ns)]: ${message}"
}

# `ensure_group` takes a group name and location
# and returns the group ID, creating it if necessary
function ensure_group() {
    local group_name=$1
    local group_location=$2

    # does group already exist?
    group_id=$(az group show --name ${group_name} \
        --output tsv --query id)

    # no, so create it
    if [[ -z "$group_id" ]]; then
        if [[ -z "$group_location" ]]; then
            # set a default location
            group_location=westus
        fi
        group_id=$(az group create \
            --name ${group_name} \
            --location ${group_location} \
            --output tsv --query id)
    fi

    echo ${group_id}
    return 0
}
