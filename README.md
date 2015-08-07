[![Build Status](https://travis-ci.org/goware/httpmock.svg?branch=master)](https://travis-ci.org/goware/httpmock)

# httpmock
Mocking 3rd party services in Go made simple


### How does it work
Httpmock runs a local server (within test) that serves predefined responses to http requests effectively faking 3rd party service in a fast, reliable way that doesn't require you to change your code, only settings (unless of course you've hardcoded service urls - that you'll have to change).

### Example
```go
package yours

import (
	"testing"
	"net/http"
	"io/ioutil"
	"github.com/goware/httpmock"
)


// the code you want to test
func Access3rdPartyService() (http.Response, error) {
	return http.Get(serviceYouDontControl)
}

// normally http://example.com/api, but we're changing it to use mock server 
// this should be the only change necessary to run tests
var serviceYouDontControl = "http://127.0.0.1:10000/api/list"

func TestSomething(t *testing.T) {

	// new mocking server
	mockService := httpmock.NewMockHTTPServer("127.0.0.1:10000")

	// define request->response pairs
	requestUrl, _ := url.Parse("http://127.0.0.1:10000/api/list")
	mockService.AddResponses([]httpmock.MockResponse{
		{
			Request: http.Request{
				Method: "GET",
				URL:    requestUrl,
			},
			Response: httpmock.Response{
				StatusCode: 200,
				Body:       "it's alive!",
			},
		},
	})

	// test code relying on 3rd party service
	serviceResponse, _ := Access3rdPartyService()
	if serviceResponse.StatusCode != 200 {
		t.Errorf("Expected status code %d, received %d", 200, serviceResponse.StatusCode)
	}
	if body, err := ioutil.ReadAll(serviceResponse.Body); err != nil {
		t.Errorf("Response body read error")
	} else {
		if string(body) != "it's alive!" {
			t.Errorf("Expected response: `%s`, received:`%s`", "it's alive!", string(body))
		}
	}
}
```
