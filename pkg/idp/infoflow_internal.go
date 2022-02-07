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

	"golang.org/x/oauth2"
)

type InfoflowInternalIdProvider struct {
	Client  *http.Client
	Config  *oauth2.Config
	AgentId string
}

func NewInfoflowInternalIdProvider(clientId string, clientSecret string, appId string, redirectUrl string) *InfoflowInternalIdProvider {
	idp := &InfoflowInternalIdProvider{}

	config := idp.getConfig(clientId, clientSecret, redirectUrl)
	idp.Config = config
	idp.AgentId = appId
	return idp
}

func (idp *InfoflowInternalIdProvider) SetHttpClient(client *http.Client) {
	idp.Client = client
}

func (idp *InfoflowInternalIdProvider) getConfig(clientId string, clientSecret string, redirectUrl string) *oauth2.Config {
	var config = &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
	}

	return config
}

type InfoflowInterToken struct {
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
}

// get more detail via: https://qy.baidu.com/doc/index.html#/inner_quickstart/flow?id=%E8%8E%B7%E5%8F%96accesstoken
func (idp *InfoflowInternalIdProvider) GetToken(code string) (*oauth2.Token, error) {
	pTokenParams := &struct {
		CorpId     string `json:"corpid"`
		Corpsecret string `json:"corpsecret"`
	}{idp.Config.ClientID, idp.Config.ClientSecret}
	resp, err := idp.Client.Get(fmt.Sprintf("https://qy.im.baidu.com/api/gettoken?corpid=%s&corpsecret=%s", pTokenParams.CorpId, pTokenParams.Corpsecret))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	pToken := &InfoflowInterToken{}
	err = json.Unmarshal(data, pToken)
	if err != nil {
		return nil, err
	}
	if pToken.Errcode != 0 {
		return nil, fmt.Errorf("pToken.Errcode = %d, pToken.Errmsg = %s", pToken.Errcode, pToken.Errmsg)
	}
	token := &oauth2.Token{
		AccessToken: pToken.AccessToken,
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
    "email": "lili4@gzdev.com",
    "imid": 40000318,
    "hiuname": "lili4",
    "status": 1,
    "extattr":
        {
            "attrs": [
                {
                    "name": "爱好",
                    "value": "旅游"
                },
                {
                    "name": "卡号,
                    "value": "1234567234"
                }
            ]
        },
   "lm": 14236463257
}
*/

type InfoflowInternalUserResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	UserId  string `json:"UserId"`
}

type InfoflowInternalUserInfo struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	UserId  string `json:"userid"`
	Imid    int    `json:"imid"`
	Name    string `json:"name"`
	Avatar  string `json:"headimg"`
	Email   string `json:"email"`
}

// get more detail via: https://qy.baidu.com/doc/index.html#/inner_serverapi/contacts?id=%e8%8e%b7%e5%8f%96%e6%88%90%e5%91%98
func (idp *InfoflowInternalIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	//Get userid first
	accessToken := token.AccessToken
	code := token.Extra("code").(string)
	resp, err := idp.Client.Get(fmt.Sprintf("https://qy.im.baidu.com/api/user/getuserinfo?access_token=%s&code=%s&agentid=%s", accessToken, code, idp.AgentId))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	userResp := &InfoflowInternalUserResp{}
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
	infoResp := &InfoflowInternalUserInfo{}
	err = json.Unmarshal(data, infoResp)
	if err != nil {
		return nil, err
	}
	if infoResp.Errcode != 0 {
		return nil, fmt.Errorf("userInfoResp.errcode = %d, userInfoResp.errmsg = %s", infoResp.Errcode, infoResp.Errmsg)
	}
	userInfo := UserInfo{
		Id:          infoResp.UserId,
		Username:    infoResp.UserId,
		DisplayName: infoResp.Name,
		AvatarUrl:   infoResp.Avatar,
		Email:       infoResp.Email,
	}

	if userInfo.Id == "" {
		userInfo.Id = userInfo.Username
	}
	return &userInfo, nil
}
