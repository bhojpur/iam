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
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	Id          string
	Username    string
	DisplayName string
	Email       string
	AvatarUrl   string
}

type IdProvider interface {
	SetHttpClient(client *http.Client)
	GetToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*UserInfo, error)
}

func GetIdProvider(providerType string, clientId string, clientSecret string, redirectUrl string) IdProvider {
	if providerType == "GitHub" {
		return NewGithubIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "Google" {
		return NewGoogleIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "QQ" {
		return NewQqIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "WeChat" {
		return NewWeChatIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "Facebook" {
		return NewFacebookIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "DingTalk" {
		return NewDingTalkIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "Weibo" {
		return NewWeiBoIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "Gitee" {
		return NewGiteeIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "LinkedIn" {
		return NewLinkedInIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "WeCom" {
		return NewWeComIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "Lark" {
		return NewLarkIdProvider(clientId, clientSecret, redirectUrl)
	} else if providerType == "GitLab" {
		return NewGitlabIdProvider(clientId, clientSecret, redirectUrl)
	} else if isGothSupport(providerType) {
		return NewGothIdProvider(providerType, clientId, clientSecret, redirectUrl)
	}

	return nil
}

var gothList = []string{"Apple", "AzureAd", "Slack"}

func isGothSupport(provider string) bool {
	for _, value := range gothList {
		if strings.EqualFold(value, provider) {
			return true
		}
	}
	return false
}
