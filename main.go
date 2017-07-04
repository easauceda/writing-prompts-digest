package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/jmoiron/jsonq"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
)

var refreshToken = os.Getenv("TOKEN")
var duration = os.Getenv("DURATION")
var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")
var emailPassword = os.Getenv("EMAIL_PASSWORD")
var userAgent = "WritingPromptsDigest/0.1 by easauceda"
var client = &http.Client{}
var auth smtp.Auth

type writingPrompt struct {
	Title   string
	URL     string
	Excerpt string
	ID      string
}

type writingPromptEmail struct {
	from    string
	to      []string
	subject string
	body    string
}

func main() {
	digest := writingPromptEmail{"easauceda@gmail.com", []string{"easauceda@gmail.com"}, "New Stories for You!", ""}
	// TODO: move 44-48 to generateDigest()
	accessToken := getAccessToken()
	prompts := getWritingPrompts(accessToken, duration)
	for i, prompt := range prompts {
		prompts[i].Excerpt = getExcerpts(prompt.ID, accessToken)
	}
	digest.body = generateDigest(prompts)
	sendEmail(digest)
}

func sendEmail(digest writingPromptEmail) {
	auth = smtp.PlainAuth("", "easauceda@gmail.com", emailPassword, "smtp.gmail.com")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + digest.subject + "\n"
	msg := []byte(subject + mime + "\n" + digest.body)
	addr := "smtp.gmail.com:587"

	if err := smtp.SendMail(addr, auth, "easauceda@gmail.com", digest.to, msg); err != nil {
		panic(err)
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

	req, _ := http.NewRequest("GET", "https://oauth.reddit.com/r/writingprompts/top.json?limit=5&t=week", nil)
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

func generateDigest(topPrompts []writingPrompt) string {
	var html bytes.Buffer
	t, err := template.New("Template").ParseFiles("template.html")
	err = t.ExecuteTemplate(&html, "template.html", topPrompts)
	if err != nil {
		panic(err)
	}
	return html.String()
}
