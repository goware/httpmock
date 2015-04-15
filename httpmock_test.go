package httpmock

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestBasicResponse(t *testing.T) {
	mockServer := NewMockHTTPServer()
	u, _ := url.Parse("http://127.0.0.1:9001/54")
	badU, _ := url.Parse("http://127.0.0.1:9001/10000")
	mockServer.AddResponses([]MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    u,
			},
			Response: Response{
				StatusCode: 200,
				Body:       "it's alive!",
			},
		},
	})

	checks := []struct {
		in  http.Request
		out Response
	}{
		{
			in: http.Request{
				Method: "GET",
				URL:    u,
			},
			out: Response{
				StatusCode: 200,
				Body:       "it's alive!",
			},
		},
		{
			in: http.Request{
				Method: "POST",
				URL:    u,
			},
			out: Response{
				StatusCode: 404,
				Body:       "route not mocked",
			},
		},
		{
			in: http.Request{
				Method: "GET",
				URL:    badU,
			},
			out: Response{
				StatusCode: 404,
				Body:       "route not mocked",
			},
		},
	}

	client := &http.Client{}

	for _, tt := range checks {
		resp, err := client.Do(&tt.in)
		if err != nil {
			t.FailNow()
		}
		if resp.StatusCode != tt.out.StatusCode {
			t.Errorf("Expected status code %d, received %d", tt.out.StatusCode, resp.StatusCode)
		}

		if resp.Body != nil {
			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Errorf("Response body read error")
			} else {
				if string(body) != tt.out.Body {
					t.Errorf("Expected response: `%s`, received:`%s`", tt.out.Body, string(body))
				}
			}
		}
	}
}
