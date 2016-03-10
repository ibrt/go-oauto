package api

import (
	"net/http"
	"github.com/ibrt/go-oauto/oauto/config"
	"encoding/json"
	"fmt"
)

type ApiHandler func(config *config.Config, r *http.Request, baseURL string) (interface{}, error)

func RegisterApiRoutes(config *config.Config) {
	http.HandleFunc("/api/authenticate", MakeHandlerFunc(config, HandleAuthenticate))
}

func MakeHandlerFunc(config *config.Config, apiHandler ApiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, fmt.Sprintf("Method '%v' is not acceptable.", r.Method), http.StatusMethodNotAllowed)
			return
		}
		if resp, err := apiHandler(config, r, fmt.Sprintf("http://%v:%v", config.RedirectHost, config.ServerPort)); err == nil {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
