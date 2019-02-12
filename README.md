# Vanpool Manager

Vanpool Manager is a web application for managing registered vanpool riders for
each day and direction. The frontend is a Vue.js app in the `web` folder. The
backend is a Go API. Data is persisted in a PostgreSQL database. Session state
is maintained in a cookie but will be moved to Redis soon. Authentication is
handled by Azure AD.

Scripts for build and deploy are in the `scripts` folder. See `.env.tpl` to
see and set needed environment variables. In Azure App Service you'll need to set
these as App Settings.

The `scripts/deploy.sh` script creates ACR, Postgres, and Web App resources in
Azure and pushes the built container from this repo to the Web App.

See [github.com/joshgav/azure-dapp][] for guidance on how to deploy this to
Kubernetes or AKS.

## License

MIT, see [LICENSE](./LICENSE).

[github.com/joshgav/azure-dapp]: https://github.com/joshgav/azure-dapp
