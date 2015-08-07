package httpmock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

func (m *MockHTTPServer) AddResponse(resp MockResponse) error {
	requestString, err := request2string(resp.Request)
	if err != nil {
		return err
	}

	m.ResponseMap[requestString] = resp.Response

	return nil
}

func (m *MockHTTPServer) AddResponses(resp []MockResponse) error {
	for _, r := range resp {
		if err := m.AddResponse(r); err != nil {
			return err
		}
	}
	return nil
}

func request2string(req http.Request) (string, error) {
	addRequestDefaults(&req)
	fragments := []string{
		req.Method,
		req.URL.RequestURI(),
	}
	if req.Body != nil {
		if body, err := ioutil.ReadAll(req.Body); err != nil {
			return "", err
		} else {
			fragments = append(fragments, string(body))
		}
	} else {
		fragments = append(fragments, "")
	}

	headerStrings := make([]string, 0)
	for index, values := range req.Header {
		headerStrings = append(headerStrings, fmt.Sprintf("%s: %s", strings.ToLower(index), strings.Join(values, ",")))
	}
	sort.Strings(headerStrings)
	fragments = append(fragments, headerStrings...)

	return strings.Join(fragments, "|"), nil
}

func addRequestDefaults(req *http.Request) {
	if req.Header == nil {
		req.Header = http.Header{}
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", defaultUserAgent)
	}

	if req.Header.Get("Accept-Encoding") == "" {
		req.Header.Set("Accept-Encoding", "gzip")
	}
}
