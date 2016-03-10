package redirect

import (
	"net/http"
	"github.com/ibrt/go-oauto/oauto/providers"
	"html/template"
	"fmt"
)

var redirectTemplate = template.Must(template.New("redirect").Parse(`
  <html>
    <body>
      <div>Token: <div id="token">{{ . }}</div></div>
     </body>
  </html>
`))

func RegisterRedirectRoutes() {
	for _, provider := range providers.Providers {
		http.HandleFunc(MakePath(provider), MakeHandlerFunc(provider))
	}
}

func MakePath(provider providers.Provider) string {
	return fmt.Sprintf("/redirect/%v", provider.GetName())
}

func MakeHandlerFunc(provider providers.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if token, err := provider.HandleRedirect(r); err == nil {
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			redirectTemplate.Execute(w, token)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}