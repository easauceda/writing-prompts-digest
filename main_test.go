package main

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func Test_GetAccessTokenResponse(t *testing.T) {
	var expectedValue = "mockAccessToken"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://www.reddit.com/api/v1/access_token",
		httpmock.NewStringResponder(200, `{"access_token":"mockAccessToken"}`))

	accessToken := getAccessToken("mockrefreshToken", "mockClientID", "mockClientSecret")
	if accessToken != expectedValue {
		t.Error("getAccessToken returned \"" + accessToken + "\" instead of \"" + expectedValue + "\"")
	}
}
