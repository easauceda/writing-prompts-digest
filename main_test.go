package main

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func Test_GetAccessTokenErrors(t *testing.T) {

	// Test Error Handling
	var expectedValue string

	expectedValue = "missing TOKEN environment variable"
	_, err := getAccessToken("", "mockClientID", "mockClientSecret")
	if err.Error() != expectedValue {
		t.Error("getAccessToken returned \"" + err.Error() + "\" instead of \"" + expectedValue + "\"")
	}

	expectedValue = "missing CLIENT_ID environment variable"
	_, err = getAccessToken("mockrefreshToken", "", "mockClientSecret")
	if err.Error() != expectedValue {
		t.Error("getAccessToken returned \"" + err.Error() + "\" instead of \"" + expectedValue + "\"")
	}

	expectedValue = "missing CLIENT_SECRET environment variable"
	_, err = getAccessToken("mockrefreshToken", "mockClientID", "")
	if err.Error() != expectedValue {
		t.Error("getAccessToken returned \"" + err.Error() + "\" instead of \"" + expectedValue + "\"")
	}
}

func Test_GetAccessTokenResponse(t *testing.T) {
	var expectedValue = "mockAccessToken"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://www.reddit.com/api/v1/access_token",
		httpmock.NewStringResponder(200, `{"access_token":"mockAccessToken"}`))

	accessToken, _ := getAccessToken("mockrefreshToken", "mockClientID", "mockClientSecret")
	if accessToken != expectedValue {
		t.Error("getAccessToken returned \"" + accessToken + "\" instead of \"" + expectedValue + "\"")
	}

	httpmock.RegisterResponder("POST", "https://www.reddit.com/api/v1/access_token",
		httpmock.NewStringResponder(404, `{"access_token":"mockAccessToken"}`))

	expectedValue = string("error requesting access token, status code 404")
	_, err := getAccessToken("mockrefreshToken", "mockClientID", "mockClientSecret")
	if err.Error() != expectedValue {
		t.Error("getAccessToken returned \"" + err.Error() + "\" instead of \"" + expectedValue + "\"")
	}

}
