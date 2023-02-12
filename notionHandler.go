package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NotionAuthConfig struct {
	Grant_type   string `json:"grant_type"`
	Code         string `json:"code"`
	Redirect_uri string `json:"redirect_uri"`
}

func formNotionAuthConfig(authCode string) *bytes.Buffer {
	body := NotionAuthConfig{
		Grant_type:   "authorization_code",
		Code:         authCode,
		Redirect_uri: getEnvValue("notionRedirectURI"),
	}
	josnfyBody, _ := json.Marshal(body)
	bufferBody := bytes.NewBuffer(josnfyBody)
	return bufferBody
}

func formNotionAuthHeader() string {
	credentials := getEnvValue("notionClientID") + ":" + getEnvValue("notionClientSecret")
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	authHeader := `Basic "` + encodedCredentials + `"`
	return authHeader
}

type NotionAuthResponse struct {
	Access_token           string      `json:"access_token"`
	Token_type             string      `json:"token_type"`
	Bot_id                 string      `json:"bot_id"`
	Workspace_name         string      `json:"workspace_name"`
	Workspace_icon         string      `json:"workspace_icon"`
	Workspace_id           string      `json:"workspace_id"`
	Owner                  NotionOwner `json:"owner"`
	Duplicated_template_id interface{} `json:"duplicated_template_id"`
}

type NotionOwner struct {
	Type string     `json:"type"`
	User NotionUser `json:"user"`
}

type NotionUser struct {
	Object     string       `json:"object"`
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Avatar_url string       `json:"avatar_url"`
	Type       string       `json:"type"`
	Person     NotionPerson `json:"person"`
	Bots       interface{}
}

type NotionPerson struct {
	Email string `json:"email"`
}

func requestNotionAuth(authCode string) NotionAuthResponse {
	requestBody := formNotionAuthConfig(authCode)
	request, err := http.NewRequest("POST", "https://api.notion.com/v1/oauth/token", requestBody)
	if err != nil {
		fmt.Printf("[ERROR] can't initial a request to Notion: %v\n", err)
	}
	request.Header.Set("Authorization", formNotionAuthHeader())
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("[ERROR] can't POST a request to Notion: %v\n", err)
	}
	defer response.Body.Close()

	d, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("[ERROR] can't read response body: %v\n", err)
	}

	var result NotionAuthResponse
	json.Unmarshal([]byte(d), &result)
	fmt.Printf("[NOTE] response Body: \n%v\n", result.Owner)
	return result
}

func showNotionUserInfo(authResponse NotionAuthResponse) {
	fmt.Printf(
		"[NOTE] user information:\nEmail: %s, \nName: %s, \nPhoto: %s\n",
		authResponse.Owner.User.Person.Email,
		authResponse.Owner.User.Name,
		authResponse.Owner.User.Avatar_url,
	)
}

func getNotionUserInfo(authResponse NotionAuthResponse) Member {
	userInfo := Member{
		Name:   authResponse.Owner.User.Name,
		Email:  authResponse.Owner.User.Person.Email,
		Source: "Notion",
		Photo:  authResponse.Owner.User.Avatar_url,
	}
	return userInfo
}
