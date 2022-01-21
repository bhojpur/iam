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
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/bhojpur/iam/pkg/utils"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/azuread"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/gitea"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/kakao"
	"github.com/markbates/goth/providers/line"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/shopify"
	"github.com/markbates/goth/providers/slack"
	"github.com/markbates/goth/providers/tumblr"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/yahoo"
	"github.com/markbates/goth/providers/yandex"
	"github.com/markbates/goth/providers/zoom"
	"golang.org/x/oauth2"
)

type GothIdProvider struct {
	Provider goth.Provider
	Session  goth.Session
}

func NewGothIdProvider(providerType string, clientId string, clientSecret string, redirectUrl string) *GothIdProvider {
	var idp GothIdProvider
	switch providerType {
	case "Amazon":
		idp = GothIdProvider{
			Provider: amazon.New(clientId, clientSecret, redirectUrl),
			Session:  &amazon.Session{},
		}
	case "Apple":
		idp = GothIdProvider{
			Provider: apple.New(clientId, clientSecret, redirectUrl, nil),
			Session:  &apple.Session{},
		}
	case "AzureAD":
		idp = GothIdProvider{
			Provider: azuread.New(clientId, clientSecret, redirectUrl, nil),
			Session:  &azuread.Session{},
		}
	case "Bitbucket":
		idp = GothIdProvider{
			Provider: bitbucket.New(clientId, clientSecret, redirectUrl),
			Session:  &bitbucket.Session{},
		}
	case "DigitalOcean":
		idp = GothIdProvider{
			Provider: digitalocean.New(clientId, clientSecret, redirectUrl),
			Session:  &digitalocean.Session{},
		}
	case "Discord":
		idp = GothIdProvider{
			Provider: discord.New(clientId, clientSecret, redirectUrl),
			Session:  &discord.Session{},
		}
	case "Dropbox":
		idp = GothIdProvider{
			Provider: dropbox.New(clientId, clientSecret, redirectUrl),
			Session:  &dropbox.Session{},
		}
	case "Facebook":
		idp = GothIdProvider{
			Provider: facebook.New(clientId, clientSecret, redirectUrl),
			Session:  &facebook.Session{},
		}
	case "Gitea":
		idp = GothIdProvider{
			Provider: gitea.New(clientId, clientSecret, redirectUrl),
			Session:  &gitea.Session{},
		}
	case "GitHub":
		idp = GothIdProvider{
			Provider: github.New(clientId, clientSecret, redirectUrl),
			Session:  &github.Session{},
		}
	case "GitLab":
		idp = GothIdProvider{
			Provider: gitlab.New(clientId, clientSecret, redirectUrl),
			Session:  &gitlab.Session{},
		}
	case "Google":
		idp = GothIdProvider{
			Provider: google.New(clientId, clientSecret, redirectUrl),
			Session:  &google.Session{},
		}
	case "Heroku":
		idp = GothIdProvider{
			Provider: heroku.New(clientId, clientSecret, redirectUrl),
			Session:  &heroku.Session{},
		}
	case "Instagram":
		idp = GothIdProvider{
			Provider: instagram.New(clientId, clientSecret, redirectUrl),
			Session:  &instagram.Session{},
		}
	case "Kakao":
		idp = GothIdProvider{
			Provider: kakao.New(clientId, clientSecret, redirectUrl),
			Session:  &kakao.Session{},
		}
	case "Linkedin":
		idp = GothIdProvider{
			Provider: linkedin.New(clientId, clientSecret, redirectUrl),
			Session:  &linkedin.Session{},
		}
	case "Line":
		idp = GothIdProvider{
			Provider: line.New(clientId, clientSecret, redirectUrl),
			Session:  &line.Session{},
		}
	case "MicrosoftOnline":
		idp = GothIdProvider{
			Provider: microsoftonline.New(clientId, clientSecret, redirectUrl),
			Session:  &microsoftonline.Session{},
		}
	case "Paypal":
		idp = GothIdProvider{
			Provider: paypal.New(clientId, clientSecret, redirectUrl),
			Session:  &paypal.Session{},
		}
	case "SalesForce":
		idp = GothIdProvider{
			Provider: salesforce.New(clientId, clientSecret, redirectUrl),
			Session:  &salesforce.Session{},
		}
	case "Shopify":
		idp = GothIdProvider{
			Provider: shopify.New(clientId, clientSecret, redirectUrl),
			Session:  &shopify.Session{},
		}
	case "Slack":
		idp = GothIdProvider{
			Provider: slack.New(clientId, clientSecret, redirectUrl),
			Session:  &slack.Session{},
		}
	case "Tumblr":
		idp = GothIdProvider{
			Provider: tumblr.New(clientId, clientSecret, redirectUrl),
			Session:  &tumblr.Session{},
		}
	case "Twitter":
		idp = GothIdProvider{
			Provider: twitter.New(clientId, clientSecret, redirectUrl),
			Session:  &twitter.Session{},
		}
	case "Yahoo":
		idp = GothIdProvider{
			Provider: yahoo.New(clientId, clientSecret, redirectUrl),
			Session:  &yahoo.Session{},
		}
	case "Yandex":
		idp = GothIdProvider{
			Provider: yandex.New(clientId, clientSecret, redirectUrl),
			Session:  &yandex.Session{},
		}
	case "Zoom":
		idp = GothIdProvider{
			Provider: zoom.New(clientId, clientSecret, redirectUrl),
			Session:  &zoom.Session{},
		}
	}

	return &idp
}

//Goth's idp all implement the Client method, but since the goth.Provider interface does not provide to modify idp's client method, reflection is required
func (idp *GothIdProvider) SetHttpClient(client *http.Client) {
	idpClient := reflect.ValueOf(idp.Provider).Elem().FieldByName("HTTPClient")
	idpClient.Set(reflect.ValueOf(client))
}

func (idp *GothIdProvider) GetToken(code string) (*oauth2.Token, error) {
	var expireAt time.Time
	//Need to construct variables supported by goth
	//to call the function to obtain accessToken
	value := url.Values{}
	value.Add("code", code)
	accessToken, err := idp.Session.Authorize(idp.Provider, value)
	//Get ExpiresAt's value
	valueOfExpire := reflect.ValueOf(idp.Session).Elem().FieldByName("ExpiresAt")
	if valueOfExpire.IsValid() {
		expireAt = valueOfExpire.Interface().(time.Time)
	}
	token := oauth2.Token{
		AccessToken: accessToken,
		Expiry:      expireAt,
	}
	return &token, err
}

func (idp *GothIdProvider) GetUserInfo(token *oauth2.Token) (*UserInfo, error) {
	gothUser, err := idp.Provider.FetchUser(idp.Session)
	if err != nil {
		return nil, err
	}
	return getUser(gothUser), nil
}

func getUser(gothUser goth.User) *UserInfo {
	user := UserInfo{
		Id:          gothUser.UserID,
		Username:    gothUser.Name,
		DisplayName: gothUser.NickName,
		Email:       gothUser.Email,
		AvatarUrl:   gothUser.AvatarURL,
	}
	//Some idp return an empty Name
	//so construct the Name with firstname and lastname or nickname
	if user.Username == "" {
		if gothUser.FirstName != "" && gothUser.LastName != "" {
			user.Username = getName(gothUser.FirstName, gothUser.LastName)
		} else {
			user.Username = gothUser.NickName
		}
	}
	if user.DisplayName == "" {
		if gothUser.FirstName != "" && gothUser.LastName != "" {
			user.DisplayName = getName(gothUser.FirstName, gothUser.LastName)
		} else {
			user.DisplayName = user.Username
		}
	}

	return &user
}

func getName(firstName, lastName string) string {
	if utils.IsChinese(firstName) || utils.IsChinese(lastName) {
		return fmt.Sprintf("%s%s", lastName, firstName)
	} else {
		return fmt.Sprintf("%s %s", firstName, lastName)
	}
}
