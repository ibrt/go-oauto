//go:generate ../../docker.sh

package main

import (
	"fmt"
	"github.com/ibrt/go-oauto/oauto/api"
	"github.com/ibrt/go-oauto/oauto/config"
	"github.com/ibrt/go-oauto/oauto/redirect"
	"github.com/kelseyhightower/envconfig"
	"net/http"
)

func main() {
	config := &config.Config{}
	envconfig.MustProcess("OAUTO", config)

	redirect.RegisterRedirectRoutes()
	api.RegisterApiRoutes(config)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", config.ServerPort), nil); err != nil {
		panic(err)
	}
}
