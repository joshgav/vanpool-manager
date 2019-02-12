# postgres for persisting registered riders
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_HOSTNAME=
# Azure creates a DB named postgres
POSTGRES_DB=postgres
POSTGRES_PORT=5432
# set to `require` for Azure
POSTGRES_SSLMODE=require

# redis for maintaining sessions
REDIS_HOSTNAME=
REDIS_PORT=6379
REDIS_PASSWORD=

# converged AAD client for identification
OAUTH_CLIENT_ID=
OAUTH_CLIENT_SECRET=
OAUTH_LOGIN_HOSTNAME=
OAUTH_LOGIN_SCHEME=https

# session key for signing cookies
SESSION_KEY=my-super-duper-session-key
