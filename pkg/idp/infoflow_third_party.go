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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type InfoflowIdProvider struct {
	Client  *http.Client
	Config  *oauth2.Config
	AgentId string
	Ticket  string
}

func NewInfoflowIdProvider(clientId string, clientSecret string, appId string, redirectUrl string) *InfoflowIdProvider {
	idp := &InfoflowIdProvider{}

	config := idp.getConfig(clientId, clientSecret, redirectUrl)
	idp.Config = config
	idp.AgentId = appId
	return idp
}

func (idp *InfoflowIdProvider) SetHttpClient(client *http.Client) {
	idp.Client = client
}

func (idp *InfoflowIdProvider) getConfig(clientId string, clientSecret string, redirectUrl string) *oauth2.Config {
	var config = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
	}

	return config
}

type InfoflowToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"suite_access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// get more detail via: https://qy.baidu.com/doc/index.html#/third_serverapi/authority
func (idp *InfoflowIdProvider) GetToken(code string) (*oauth2.Token, error) {
	pTokenParams := &struct {
		SuiteId     string `json:"suite_id"`
		SuiteSecret string `json:"suite_secret"`
		SuiteTicket string `json:"suite_ticket"`
	}{idp.Config.ClientID, idp.Config.ClientSecret, idp.Ticket}
	data, err := idp.postWithBody(pTokenParams, "https://api.im.baidu.com/api/service/get_suite_token")

	pToken := &InfoflowToken{}
	err = json.Unmarshal(data, pToken)
	if err != nil {
		return nil, err
	}
	if pToken.Errcode != 0 {
		return nil, fmt.Errorf("pToken.Errcode = %d, pToken.Errmsg = %s", pToken.Errcode, pToken.Errmsg)
	}
	token := &oauth2.Token{
		AccessToken: pToken.AccessToken,
		Expiry:      time.Unix(time.Now().Unix()+int64(pToken.ExpiresIn), 0),
	}

	raw := make(map[string]interface{})
	raw["code"] = code
	token = token.WithExtra(raw)

	return token, nil
}

/*
{
    "errcode": 0,
    "errmsg": "ok",
    "userid": "lili",
    "name": "丽丽",
    "department": [1],
    "mobile": "13500088888",
    "email": "pramila@bhojpur.net",
    "imid": 40000318,
    "hiuname": "lili4",
    "status": 1,
    "extattr": {
        "attrs": [
            {
                "name": "爱好",
                "value": "旅游"
            },
            {
                "name": "卡号",
                "value": "1234567234"
            }
        ]
    },
    "lm" : 14236463257
}
*/

type InfoflowUserResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	UserId  string `json:"UserId"`
}

type InfoflowUserInfo struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Imid    string `json:"imid"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

// get more detail via: https://qy.baidu.com/doc/index.html#/third_serverapi/contacts?id=%e8%8e%b7%e5%8f%96%e6%88%90%e5%91%98
func (idp *InfoflowIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	//Get userid first
	accessToken := token.AccessToken
	code := token.Extra("code").(string)
	resp, err := idp.Client.Get(fmt.Sprintf("https://api.im.baidu.com/api/user/getuserinfo?access_token=%s&code=%s&agentid=%s", accessToken, code, idp.AgentId))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	userResp := &InfoflowUserResp{}
	err = json.Unmarshal(data, userResp)
	if err != nil {
		return nil, err
	}
	if userResp.Errcode != 0 {
		return nil, fmt.Errorf("userIdResp.Errcode = %d, userIdResp.Errmsg = %s", userResp.Errcode, userResp.Errmsg)
	}
	//Use userid and accesstoken to get user information
	resp, err = idp.Client.Get(fmt.Sprintf("https://api.im.baidu.com/api/user/get?access_token=%s&userid=%s", accessToken, userResp.UserId))
	if err != nil {
		return nil, err
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	infoResp := &InfoflowUserInfo{}
	err = json.Unmarshal(data, infoResp)
	if err != nil {
		return nil, err
	}
	if infoResp.Errcode != 0 {
		return nil, fmt.Errorf("userInfoResp.errcode = %d, userInfoResp.errmsg = %s", infoResp.Errcode, infoResp.Errmsg)
	}
	userInfo := UserInfo{
		Id:          infoResp.Imid,
		Username:    infoResp.Name,
		DisplayName: infoResp.Name,
		Email:       infoResp.Email,
	}

	if userInfo.Id == "" {
		userInfo.Id = userInfo.Username
	}
	return &userInfo, nil
}

func (idp *InfoflowIdProvider) postWithBody(body interface{}, url string) ([]byte, error) {
	bs, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	r := strings.NewReader(string(bs))
	resp, err := idp.Client.Post(url, "application/json;charset=UTF-8", r)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	return data, nil
}
