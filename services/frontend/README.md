### frontend

This service serves an HTML/JS site for vanpool-manager. It is powered by the
[Buffalo][] server development framework.

`make azure-infra` sets up Azure infrastructure. `make local-infra` sets up
local container infrastructure. Other tasks are grouped by targets in `./hack/targets`.

`buffalo dev` builds and starts the containerized service.

[Buffalo]: https://github.com/gobuffalo/buffalo
