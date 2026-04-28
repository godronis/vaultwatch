package alert

import (
	"net/http"
	"strings"
)

// rewriteToServer returns an http.RoundTripper that redirects all requests
// to the given base URL, preserving the path and query. This is used in tests
// to intercept outbound HTTP calls without modifying production code.
func rewriteToServer(baseURL string) http.RoundTripper {
	return &redirectTransport{base: strings.TrimRight(baseURL, "/")}
}

type redirectTransport struct {
	base string
}

func (r *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	cloned.URL.Scheme = "http"
	cloned.URL.Host = strings.TrimPrefix(r.base, "http://")
	return http.DefaultTransport.RoundTrip(cloned)
}
