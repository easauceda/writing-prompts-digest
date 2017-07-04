package main

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/jmoiron/jsonq"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var refreshToken = os.Getenv("TOKEN")
var duration = os.Getenv("DURATION")
var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")
var userAgent = "WritingPromptsDigest/0.1 by easauceda"
var client = &http.Client{}

type writingPrompt struct {
	Title   string
	URL     string
	Excerpt string
	ID      string
}

func main() {
	accessToken := getAccessToken()
	prompts := getWritingPrompts(accessToken, duration)

	for _, prompt := range prompts {
		prompt.excerpt = getExcerpts(prompt.ID, accessToken)
	}
}

func getExcerpts(promptID string, accessToken string) string {
	req, _ := http.NewRequest("GET", "https://oauth.reddit.com/r/writingprompts/comments/"+promptID+".json?depth=1&limit=2", nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	respStr, _ := ioutil.ReadAll(resp.Body)
	testStr, _ := jsonparser.GetString(respStr, "[1]", "data", "children", "[1]", "data", "body")
	return fmt.Sprintf("%.500s...\n", testStr)
}

func getAccessToken() string {
	var tokenResp map[string]interface{}
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, _ := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	json.NewDecoder(resp.Body).Decode(&tokenResp)
	defer resp.Body.Close()
	tokenJSON := jsonq.NewQuery(tokenResp)
	token, _ := tokenJSON.String("access_token")
	return token
}

func getWritingPrompts(accessToken string, duration string) []writingPrompt {
	var writingPrompts = make([]writingPrompt, 0)
	var promptResp map[string]interface{}

	req, _ := http.NewRequest("GET", "https://oauth.reddit.com/r/writingprompts/top.json?limit=5", nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&promptResp)
	promptsJSON := jsonq.NewQuery(promptResp)
	prompts, _ := promptsJSON.ArrayOfObjects("data", "children")

	for _, promptJSON := range prompts {
		prompt := jsonq.NewQuery(promptJSON)
		promptTitle, err := prompt.String("data", "title")
		if err != nil {
			panic(err)
		}
		promptID, err := prompt.String("data", "id")
		if err != nil {
			panic(err)
		}
		promptURL, err := prompt.String("data", "url")
		if err != nil {
			panic(err)
		}
		newWritingPrompt := writingPrompt{Title: promptTitle, ID: promptID, URL: promptURL}
		writingPrompts = append(writingPrompts, newWritingPrompt)
	}
	return writingPrompts
}

func generateDigest(topPrompts []writingPrompt) {

}
