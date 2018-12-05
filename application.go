package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	slackClientID     = "YOUR_SLACK_CLIENT_ID"
	slackClientSecret = "YOUR_SLACK_CLIENT_SECRET"
	slackOAuthURL     = "https://slack.com/api/oauth.access"
)

type slackOAuthResponse struct {
	Ok              bool   `json:"ok"`
	AccessToken     string `json:"access_token"`
	Scope           string `json:"scope"`
	UserID          string `json:"user_id"`
	TeamName        string `json:"team_name"`
	TeamID          string `json:"team_id"`
	IncomingWebhook struct {
		Channel          string `json:"channel"`
		ChannelID        string `json:"channel_id"`
		ConfigurationURL string `json:"configuration_url"`
		URL              string `json:"url"`
	} `json:"incoming_webhook"`
	Bot struct {
		BotUserID      string `json:"bot_user_id"`
		BotAccessToken string `json:"bot_access_token"`
	} `json:"bot"`
}

type tmplPageData struct {
	BotToken string
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("main.html"))

	dataTmpl := tmplPageData{
		BotToken: "",
	}

	tmpl.Execute(w, dataTmpl)
}

func main() {
	http.HandleFunc("/add_gopher", slackOAuth)
	http.HandleFunc("/", mainPage)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))
	if err := http.ListenAndServe(":5000", nil); err != nil {
		panic(err)
	}
}
func slackOAuth(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")

	client := &http.Client{}

	values := url.Values{
		"client_id":     {slackClientID},
		"client_secret": {slackClientSecret},
		"code":          {code},
	}
	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", slackOAuthURL, reqBody)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}

	body, _ := ioutil.ReadAll(resp.Body)
	var data slackOAuthResponse
	json.Unmarshal(body, &data)

	fmt.Println(time.Now().Format("2006.01.02 15:04:05"), " -- Another SaltGopher added :)")

	tmpl := template.Must(template.ParseFiles("resp.html"))

	dataTmpl := tmplPageData{
		BotToken: data.Bot.BotAccessToken,
	}

	tmpl.Execute(w, dataTmpl)

}
