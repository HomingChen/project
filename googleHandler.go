package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gOAuth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

func formGoogleAuthConfig() *oauth2.Config {
	credentials := []byte(getEnvValue("googleCredentials"))
	config, err := google.ConfigFromJSON(credentials, tasks.TasksReadonlyScope, gOAuth2.UserinfoEmailScope, gOAuth2.UserinfoProfileScope)
	if err != nil {
		log.Fatalf("[ERROR] Unable to parse google credentials string to config: %v\n", err)
	}
	return config
}

func getGoogleAuthURL(config *oauth2.Config) string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("[Note] google authURL: %s\n", authURL)
	return authURL
}

func getGoogleToken(config *oauth2.Config, authCode string) *oauth2.Token {
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("[ERROR] Unable to retrieve token from web: %v\n", err)
	}
	return token
}

func requestGoogleAuth(authCode string) *http.Client {
	config := formGoogleAuthConfig()
	token := getGoogleToken(config, authCode)
	client := config.Client(context.Background(), token)
	return client
}

func initialGoogleTasksService(client *http.Client) *tasks.Service {
	ctx := context.Background()
	service, err := tasks.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("[ERROR] Unable to retrieve tasks Client %v\n", err)
	}
	return service
}

func initialGoogleOAuthService(client *http.Client) *gOAuth2.Service {
	ctx := context.Background()
	service, err := gOAuth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("[ERROR] Unable to retrieve google OAuth2 Client %v\n", err)
	}
	return service
}

func showGoogleTasksList(tasksService *tasks.Service) {
	tasksList, err := tasksService.Tasklists.List().Do()
	if err != nil {
		log.Fatalf("[ERROR] Unable to retrieve task lists: %v\n", err)
	}

	fmt.Println("[NOTE] Task Lists:")
	if len(tasksList.Items) > 0 {
		for _, i := range tasksList.Items {
			fmt.Printf("%s (%s)\n", i.Title, i.Id)
		}
	} else {
		fmt.Print("[NOTE] No task lists found.")
	}
}

func showGoogleUserInfo(gOAuth2Service gOAuth2.Service) {
	userInfo, err := gOAuth2Service.Userinfo.Get().Do()
	if err != nil {
		log.Fatalf("[ERROR] Unable to retrieve user information: %v", err)
	}
	fmt.Printf("[NOTE] user information:\nEmail: %s, \nName: %s, \nPhoto: %s\n", userInfo.Email, userInfo.Name, userInfo.Picture)
}

func getGoogleUserInfo(gOAuth2Service gOAuth2.Service) Member {
	userInfo, err := gOAuth2Service.Userinfo.Get().Do()
	if err != nil {
		log.Fatalf("[ERROR] Unable to retrieve user information: %v", err)
	}
	memberData := Member{
		Name:   userInfo.Name,
		Email:  userInfo.Email,
		Source: "google",
		Photo:  userInfo.Picture,
	}
	return memberData
}
