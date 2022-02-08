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
	"errors"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type GoogleIdProvider struct {
	Client *http.Client
	Config *oauth2.Config
}

func NewGoogleIdProvider(clientId string, clientSecret string, redirectUrl string) *GoogleIdProvider {
	idp := &GoogleIdProvider{}

	config := idp.getConfig()
	config.ClientID = clientId
	config.ClientSecret = clientSecret
	config.RedirectURL = redirectUrl
	idp.Config = config

	return idp
}

func (idp *GoogleIdProvider) SetHttpClient(client *http.Client) {
	idp.Client = client
}

func (idp *GoogleIdProvider) getConfig() *oauth2.Config {
	var endpoint = oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://accounts.google.com/o/oauth2/token",
	}

	var config = &oauth2.Config{
		Scopes:   []string{"profile", "email"},
		Endpoint: endpoint,
	}

	return config
}

func (idp *GoogleIdProvider) GetToken(code string) (*oauth2.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, idp.Client)
	return idp.Config.Exchange(ctx, code)
}

//{
//	"id": "110613473084924141234",
//	"email": "info@bhojpur.net",
//	"verified_email": true,
//	"name": "Shashi Bhushan Rai",
//	"given_name": "Shashi Bhushan",
//	"family_name": "Rai",
//	"picture": "https://static.bhojpur.net/image/avatars/photo.png",
//	"locale": "en"
//}

type GoogleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (idp *GoogleIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?alt=json&access_token=%s", token.AccessToken)
	resp, err := idp.Client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var googleUserInfo GoogleUserInfo
	err = json.Unmarshal(body, &googleUserInfo)
	if err != nil {
		return nil, err
	}

	if googleUserInfo.Email == "" {
		return nil, errors.New("google email is empty")
	}

	userInfo := UserInfo{
		Id:          googleUserInfo.Id,
		Username:    googleUserInfo.Email,
		DisplayName: googleUserInfo.Name,
		Email:       googleUserInfo.Email,
		AvatarUrl:   googleUserInfo.Picture,
	}
	return &userInfo, nil
}
