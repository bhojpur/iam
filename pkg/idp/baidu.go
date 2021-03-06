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
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type BaiduIdProvider struct {
	Client *http.Client
	Config *oauth2.Config
}

func NewBaiduIdProvider(clientId string, clientSecret string, redirectUrl string) *BaiduIdProvider {
	idp := &BaiduIdProvider{}

	config := idp.getConfig()
	config.ClientID = clientId
	config.ClientSecret = clientSecret
	config.RedirectURL = redirectUrl
	idp.Config = config

	return idp
}

func (idp *BaiduIdProvider) SetHttpClient(client *http.Client) {
	idp.Client = client
}

func (idp *BaiduIdProvider) getConfig() *oauth2.Config {
	var endpoint = oauth2.Endpoint{
		AuthURL:  "https://openapi.baidu.com/oauth/2.0/authorize",
		TokenURL: "https://openapi.baidu.com/oauth/2.0/token",
	}

	var config = &oauth2.Config{
		Scopes:   []string{"email"},
		Endpoint: endpoint,
	}

	return config
}

func (idp *BaiduIdProvider) GetToken(code string) (*oauth2.Token, error) {
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, idp.Client)
	return idp.Config.Exchange(ctx, code)
}

/*
{
    "userid":"2097322476",
    "username":"wl19871011",
    "realname":"??????",
    "userdetail":"????????????",
    "birthday":"1973-08-25",
    "marriage":"??????",
    "sex":"???",
    "blood":"A+",
    "constellation":"??????",
    "figure":"??????",
    "education":"??????/??????",
    "trade":"?????????/????????????",
    "job":"??????",
    "birthday_year":"1973",
    "birthday_month":"08",
    "birthday_day":"25",
}
*/

type BaiduUserInfo struct {
	OpenId   string `json:"openid"`
	Username string `json:"username"`
	Portrait string `json:"portrait"`
}

func (idp *BaiduIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	resp, err := idp.Client.Get(fmt.Sprintf("https://openapi.baidu.com/rest/2.0/passport/users/getInfo?access_token=%s", token.AccessToken))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	baiduUser := BaiduUserInfo{}
	if err = json.Unmarshal(data, &baiduUser); err != nil {
		return nil, err
	}

	userInfo := UserInfo{
		Id:          baiduUser.OpenId,
		Username:    baiduUser.Username,
		DisplayName: baiduUser.Username,
		AvatarUrl:   fmt.Sprintf("https://himg.bdimg.com/sys/portrait/item/%s", baiduUser.Portrait),
	}
	return &userInfo, nil
}
