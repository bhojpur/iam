package idp

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

type GithubIdProvider struct {
	Client *http.Client
	Config *oauth2.Config
}

func NewGithubIdProvider(clientId string, clientSecret string, redirectUrl string) *GithubIdProvider {
	idp := &GithubIdProvider{}

	config := idp.getConfig()
	config.ClientID = clientId
	config.ClientSecret = clientSecret
	config.RedirectURL = redirectUrl
	idp.Config = config

	return idp
}

func (idp *GithubIdProvider) SetHttpClient(client *http.Client) {
	idp.Client = client
}

func (idp *GithubIdProvider) getConfig() *oauth2.Config {
	var endpoint = oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	}

	var config = &oauth2.Config{
		Scopes:   []string{"user:email", "read:user"},
		Endpoint: endpoint,
	}

	return config
}

func (idp *GithubIdProvider) GetToken(code string) (*oauth2.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, idp.Client)
	return idp.Config.Exchange(ctx, code)
}

//{
//	"login": "shashi-rai",
//	"id": 3781234,
//	"node_id": "MDQ6VXNlcjM3O123456=",
//	"avatar_url": "https://avatars.githubusercontent.com/u/3781234?v=4",
//	"gravatar_id": "",
//	"url": "https://api.github.com/users/shashi-rai",
//	"html_url": "https://github.com/shashi-rai",
//	"followers_url": "https://api.github.com/users/shashi-rai/followers",
//	"following_url": "https://api.github.com/users/shashi-rai/following{/other_user}",
//	"gists_url": "https://api.github.com/users/shashi-rai/gists{/gist_id}",
//	"starred_url": "https://api.github.com/users/shashi-rai/starred{/owner}{/repo}",
//	"subscriptions_url": "https://api.github.com/users/shashi-rai/subscriptions",
//	"organizations_url": "https://api.github.com/users/shashi-rai/orgs",
//	"repos_url": "https://api.github.com/users/shashi-rai/repos",
//	"events_url": "https://api.github.com/users/shashi-rai/events{/privacy}",
//	"received_events_url": "https://api.github.com/users/shashi-rai/received_events",
//	"type": "User",
//	"site_admin": false,
//	"name": "Shashi Bhushan Rai",
//	"company": "Bhojpur Consulting",
//	"blog": "https://blog.bhojpur-consulting.com",
//	"location": "Arrah, Bihar, India",
//	"email": "info@bhojpur.net",
//	"hireable": true,
//	"bio": "My bio",
//	"twitter_username": null,
//	"public_repos": 45,
//	"public_gists": 3,
//	"followers": 123,
//	"following": 31,
//	"created_at": "2016-03-06T13:16:13Z",
//	"updated_at": "2020-05-30T12:15:29Z",
//	"private_gists": 0,
//	"total_private_repos": 12,
//	"owned_private_repos": 12,
//	"disk_usage": 46331,
//	"collaborators": 5,
//	"two_factor_authentication": true,
//	"plan": {
//		"name": "free",
//		"space": 976562499,
//		"collaborators": 0,
//		"private_repos": 10000
//	}
//}

type GitHubUserInfo struct {
	Login                   string      `json:"login"`
	Id                      int         `json:"id"`
	NodeId                  string      `json:"node_id"`
	AvatarUrl               string      `json:"avatar_url"`
	GravatarId              string      `json:"gravatar_id"`
	Url                     string      `json:"url"`
	HtmlUrl                 string      `json:"html_url"`
	FollowersUrl            string      `json:"followers_url"`
	FollowingUrl            string      `json:"following_url"`
	GistsUrl                string      `json:"gists_url"`
	StarredUrl              string      `json:"starred_url"`
	SubscriptionsUrl        string      `json:"subscriptions_url"`
	OrganizationsUrl        string      `json:"organizations_url"`
	ReposUrl                string      `json:"repos_url"`
	EventsUrl               string      `json:"events_url"`
	ReceivedEventsUrl       string      `json:"received_events_url"`
	Type                    string      `json:"type"`
	SiteAdmin               bool        `json:"site_admin"`
	Name                    string      `json:"name"`
	Company                 string      `json:"company"`
	Blog                    string      `json:"blog"`
	Location                string      `json:"location"`
	Email                   string      `json:"email"`
	Hireable                bool        `json:"hireable"`
	Bio                     string      `json:"bio"`
	TwitterUsername         interface{} `json:"twitter_username"`
	PublicRepos             int         `json:"public_repos"`
	PublicGists             int         `json:"public_gists"`
	Followers               int         `json:"followers"`
	Following               int         `json:"following"`
	CreatedAt               time.Time   `json:"created_at"`
	UpdatedAt               time.Time   `json:"updated_at"`
	PrivateGists            int         `json:"private_gists"`
	TotalPrivateRepos       int         `json:"total_private_repos"`
	OwnedPrivateRepos       int         `json:"owned_private_repos"`
	DiskUsage               int         `json:"disk_usage"`
	Collaborators           int         `json:"collaborators"`
	TwoFactorAuthentication bool        `json:"two_factor_authentication"`
	Plan                    struct {
		Name          string `json:"name"`
		Space         int    `json:"space"`
		Collaborators int    `json:"collaborators"`
		PrivateRepos  int    `json:"private_repos"`
	} `json:"plan"`
}

func (idp *GithubIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "token "+token.AccessToken)
	resp, err := idp.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var githubUserInfo GitHubUserInfo
	err = json.Unmarshal(body, &githubUserInfo)
	if err != nil {
		return nil, err
	}

	userInfo := UserInfo{
		Id:          strconv.Itoa(githubUserInfo.Id),
		Username:    githubUserInfo.Login,
		DisplayName: githubUserInfo.Name,
		Email:       githubUserInfo.Email,
		AvatarUrl:   githubUserInfo.AvatarUrl,
	}
	return &userInfo, nil
}
