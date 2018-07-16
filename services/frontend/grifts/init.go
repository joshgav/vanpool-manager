package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/joshgav/vanpool-manager/services/frontend/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
