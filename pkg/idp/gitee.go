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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type GiteeIdProvider struct {
	Client *http.Client
	Config *oauth2.Config
}

func NewGiteeIdProvider(clientId string, clientSecret string, redirectUrl string) *GiteeIdProvider {
	idp := &GiteeIdProvider{}

	config := idp.getConfig(clientId, clientSecret, redirectUrl)
	idp.Config = config

	return idp
}

func (idp *GiteeIdProvider) SetHttpClient(client *http.Client) {
	idp.Client = client
}

// getConfig return a point of Config, which describes a typical 3-legged OAuth2 flow
func (idp *GiteeIdProvider) getConfig(clientId string, clientSecret string, redirectUrl string) *oauth2.Config {
	var endpoint = oauth2.Endpoint{
		TokenURL: "https://gitee.com/oauth/token",
	}

	var config = &oauth2.Config{
		Scopes: []string{"user_info emails"},

		Endpoint:     endpoint,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
	}

	return config
}

type GiteeAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	CreatedAt    int    `json:"created_at"`
}

// GetToken use code get access_token (*operation of getting code ought to be done in front)
// The POST Url format of submission is: https://gitee.com/oauth/token?grant_type=authorization_code&code={code}&client_id={client_id}&redirect_uri={redirect_uri}&client_secret={client_secret}
// get more detail via: https://gitee.com/api/v5/oauth_doc#/
func (idp *GiteeIdProvider) GetToken(code string) (*oauth2.Token, error) {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", idp.Config.ClientID)
	params.Add("client_secret", idp.Config.ClientSecret)
	params.Add("code", code)
	params.Add("redirect_uri", idp.Config.RedirectURL)

	accessTokenUrl := fmt.Sprintf("%s?%s", idp.Config.Endpoint.TokenURL, params.Encode())
	bs, _ := json.Marshal(params.Encode())
	req, _ := http.NewRequest("POST", accessTokenUrl, strings.NewReader(string(bs)))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.101 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	rbs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	tokenResp := GiteeAccessToken{}
	if err = json.Unmarshal(rbs, &tokenResp); err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       time.Unix(time.Now().Unix()+int64(tokenResp.ExpiresIn), 0),
	}

	return token, nil
}

/*
{
    "id": 999999,
    "login": "xxx",
    "name": "xxx",
    "avatar_url": "https://static.bhojpur.net/image/avatars/no_portrait.png",
    "url": "https://gitee.com/api/v5/users/xxx",
    "html_url": "https://gitee.com/xxx",
    "followers_url": "https://gitee.com/api/v5/users/xxx/followers",
    "following_url": "https://gitee.com/api/v5/users/xxx/following_url{/other_user}",
    "gists_url": "https://gitee.com/api/v5/users/xxx/gists{/gist_id}",
    "starred_url": "https://gitee.com/api/v5/users/xxx/starred{/owner}{/repo}",
    "subscriptions_url": "https://gitee.com/api/v5/users/xxx/subscriptions",
    "organizations_url": "https://gitee.com/api/v5/users/xxx/orgs",
    "repos_url": "https://gitee.com/api/v5/users/xxx/repos",
    "events_url": "https://gitee.com/api/v5/users/xxx/events{/privacy}",
    "received_events_url": "https://gitee.com/api/v5/users/xxx/received_events",
    "type": "User",
    "blog": null,
    "weibo": null,
    "bio": "个人博客：https://gitee.com/xxx/xxx/pages",
    "public_repos": 2,
    "public_gists": 0,
    "followers": 0,
    "following": 0,
    "stared": 0,
    "watched": 2,
    "created_at": "2019-08-03T23:21:16+08:00",
    "updated_at": "2021-06-14T12:47:09+08:00",
    "email": null
}
*/

type GiteeUserResponse struct {
	AvatarUrl         string `json:"avatar_url"`
	Bio               string `json:"bio"`
	Blog              string `json:"blog"`
	CreatedAt         string `json:"created_at"`
	Email             string `json:"email"`
	EventsUrl         string `json:"events_url"`
	Followers         int    `json:"followers"`
	FollowersUrl      string `json:"followers_url"`
	Following         int    `json:"following"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	HtmlUrl           string `json:"html_url"`
	Id                int    `json:"id"`
	Login             string `json:"login"`
	MemberRole        string `json:"member_role"`
	Name              string `json:"name"`
	OrganizationsUrl  string `json:"organizations_url"`
	PublicGists       int    `json:"public_gists"`
	PublicRepos       int    `json:"public_repos"`
	ReceivedEventsUrl string `json:"received_events_url"`
	ReposUrl          string `json:"repos_url"`
	Stared            int    `json:"stared"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	UpdatedAt         string `json:"updated_at"`
	Url               string `json:"url"`
	Watched           int    `json:"watched"`
	Weibo             string `json:"weibo"`
}

// GetUserInfo Use userid and access_token to get UserInfo
// get more detail via: https://gitee.com/api/v5/swagger#/getV5User
func (idp *GiteeIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	var gtUserInfo GiteeUserResponse
	accessToken := token.AccessToken

	u := fmt.Sprintf("https://gitee.com/api/v5/user?access_token=%s",
		accessToken)

	userinfoResp, err := idp.GetUrlResp(u)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(userinfoResp), &gtUserInfo); err != nil {
		return nil, err
	}

	userInfo := UserInfo{
		Id:          strconv.Itoa(gtUserInfo.Id),
		Username:    gtUserInfo.Name,
		DisplayName: gtUserInfo.Name,
		Email:       gtUserInfo.Email,
		AvatarUrl:   gtUserInfo.AvatarUrl,
	}

	return &userInfo, nil
}

func (idp *GiteeIdProvider) GetUrlResp(url string) (string, error) {
	resp, err := idp.Client.Get(url)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
