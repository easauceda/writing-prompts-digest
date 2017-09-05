package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"

	"strconv"

	"github.com/buger/jsonparser"
	"github.com/jmoiron/jsonq"
	log "github.com/sirupsen/logrus"
)

var refreshToken = os.Getenv("TOKEN")
var duration = os.Getenv("DURATION")
var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")
var emailPassword = os.Getenv("EMAIL_PASSWORD")
var emailAddress = os.Getenv("EMAIL_ADDRESS")
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
	log.Info("Generating Writing Prompts Digest for ", time.Now().Local().Format("Mon Jan 1, 2006"))

	accessToken := getAccessToken(refreshToken, clientID, clientSecret)
	prompts := getWritingPrompts(&accessToken)
	digestBody := generateDigest(prompts)

	digest := writingPromptEmail{emailAddress, []string{emailAddress}, "New Stories for You!", digestBody}
	sendEmail(digest)
}

func sendEmail(digest writingPromptEmail) {
	auth = smtp.PlainAuth("", emailAddress, emailPassword, "smtp.gmail.com")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + digest.subject + "\n"
	msg := []byte(subject + mime + "\n" + digest.body)
	addr := "smtp.gmail.com:587"

	err := smtp.SendMail(addr, auth, emailAddress, digest.to, msg)
	if err != nil {
		panic(err)
	}
}

func getExcerpts(promptID string, accessToken *string) string {
	req, _ := http.NewRequest("GET", "https://oauth.reddit.com/r/writingprompts/comments/"+promptID+".json?depth=1&limit=2", nil)
	req.Header.Add("Authorization", "Bearer "+*accessToken)
	req.Header.Set("User-Agent", userAgent)

	log.Debug("Requesting Excerpt for ", promptID)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error getting an except from ", promptID, err)
	}
	respStr, _ := ioutil.ReadAll(resp.Body)
	testStr, _ := jsonparser.GetString(respStr, "[1]", "data", "children", "[1]", "data", "body")
	return fmt.Sprintf("%.500s...\n", testStr)
}

func getAccessToken(refreshToken string, clientID string, clientSecret string) string {
	var tokenResp map[string]interface{}

	reqParams := url.Values{}
	reqParams.Set("grant_type", "refresh_token")
	reqParams.Set("refresh_token", refreshToken)

	req, _ := http.NewRequest("POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(reqParams.Encode()))
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		log.Fatal("error requesting access token, status code " + strconv.Itoa(resp.StatusCode))
	}
	if err != nil {
		log.Fatal("Error requesting Access Token", err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&tokenResp)
	tokenJSON := jsonq.NewQuery(tokenResp)
	token, err := tokenJSON.String("access_token")
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func getWritingPrompts(accessToken *string) []writingPrompt {
	var writingPrompts = make([]writingPrompt, 0)
	var promptResp map[string]interface{}

	req, _ := http.NewRequest("GET", "https://oauth.reddit.com/r/writingprompts/top.json?limit=5&t=day", nil)
	req.Header.Add("Authorization", "Bearer "+*accessToken)
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
		promptExcerpt := getExcerpts(promptID, accessToken)
		newWritingPrompt := writingPrompt{Title: promptTitle, ID: promptID, URL: promptURL, Excerpt: promptExcerpt}
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
