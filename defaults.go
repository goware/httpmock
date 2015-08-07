// +build !go1.5

package httpmock

import "net/http"

func addRequestDefaults(req *http.Request) {
	if req.Header == nil {
		req.Header = http.Header{}
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Go 1.1 package http")
	}

	if req.Header.Get("Accept-Encoding") == "" {
		req.Header.Set("Accept-Encoding", "gzip")
	}
}
