//go:generate ../../docker.sh

package main

import (
	"net/http"
	"github.com/ibrt/go-oauto/oauto/redirect"
	"github.com/ibrt/go-oauto/oauto/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/ibrt/go-oauto/oauto/api"
	"fmt"
)

func main() {
	config := &config.Config{}
	envconfig.MustProcess("OAUTO", config)

	redirect.RegisterRedirectRoutes()
	api.RegisterApiRoutes(config)

	http.ListenAndServe(fmt.Sprintf(":%v", config.ServerPort), nil)
}
