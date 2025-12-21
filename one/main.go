package one

import (
	"net/http"

	"github.com/timefactoryio/frame/fx"
)

type One struct {
	fx.Fx
	mux *http.ServeMux
}

func NewOne(pathlessUrl, apiURL string) *One {
	f := &One{
		Fx: fx.NewFx(pathlessUrl, apiURL),
	}
	return f
}
