package object

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
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/bhojpur/iam/pkg/utils"
	"xorm.io/core"
)

type Code struct {
	Message string `xorm:"varchar(100)" json:"message"`
	Code    string `xorm:"varchar(100)" json:"code"`
}

type Token struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Application  string `xorm:"varchar(100)" json:"application"`
	Organization string `xorm:"varchar(100)" json:"organization"`
	User         string `xorm:"varchar(100)" json:"user"`

	Code          string `xorm:"varchar(100)" json:"code"`
	AccessToken   string `xorm:"mediumtext" json:"accessToken"`
	RefreshToken  string `xorm:"mediumtext" json:"refreshToken"`
	ExpiresIn     int    `json:"expiresIn"`
	Scope         string `xorm:"varchar(100)" json:"scope"`
	TokenType     string `xorm:"varchar(100)" json:"tokenType"`
	CodeChallenge string `xorm:"varchar(100)" json:"codeChallenge"`
}

type TokenWrapper struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func GetTokenCount(owner, field, value string) int {
	session := adapter.Engine.Where("owner=?", owner)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Token{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetTokens(owner string) []*Token {
	tokens := []*Token{}
	err := adapter.Engine.Desc("created_time").Find(&tokens, &Token{Owner: owner})
	if err != nil {
		panic(err)
	}

	return tokens
}

func GetPaginationTokens(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Token {
	tokens := []*Token{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&tokens)
	if err != nil {
		panic(err)
	}

	return tokens
}

func getToken(owner string, name string) *Token {
	if owner == "" || name == "" {
		return nil
	}

	token := Token{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&token)
	if err != nil {
		panic(err)
	}

	if existed {
		return &token
	}

	return nil
}

func getTokenByCode(code string) *Token {
	token := Token{}
	existed, err := adapter.Engine.Where("code=?", code).Get(&token)
	if err != nil {
		panic(err)
	}

	if existed {
		return &token
	}

	return nil
}

func GetToken(id string) *Token {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getToken(owner, name)
}

func UpdateToken(id string, token *Token) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getToken(owner, name) == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(token)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddToken(token *Token) bool {
	affected, err := adapter.Engine.Insert(token)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteToken(token *Token) bool {
	affected, err := adapter.Engine.ID(core.PK{token.Owner, token.Name}).Delete(&Token{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func CheckOAuthLogin(clientId string, responseType string, redirectUri string, scope string, state string) (string, *Application) {
	if responseType != "code" {
		return "response_type should be \"code\"", nil
	}

	application := GetApplicationByClientId(clientId)
	if application == nil {
		return "Invalid client_id", nil
	}

	validUri := false
	for _, tmpUri := range application.RedirectUris {
		if strings.Contains(redirectUri, tmpUri) {
			validUri = true
			break
		}
	}
	if !validUri {
		return fmt.Sprintf("Redirect URI: \"%s\" doesn't exist in the allowed Redirect URI list", redirectUri), application
	}

	// Mask application for /api/get-app-login
	application.ClientSecret = ""
	return "", application
}

func GetOAuthCode(userId string, clientId string, responseType string, redirectUri string, scope string, state string, nonce string, challenge string) *Code {
	user := GetUser(userId)
	if user == nil {
		return &Code{
			Message: fmt.Sprintf("The user: %s doesn't exist", userId),
			Code:    "",
		}
	}
	if user.IsForbidden {
		return &Code{
			Message: "error: the user is forbidden to sign in, please contact the administrator",
			Code:    "",
		}
	}

	msg, application := CheckOAuthLogin(clientId, responseType, redirectUri, scope, state)
	if msg != "" {
		return &Code{
			Message: msg,
			Code:    "",
		}
	}

	accessToken, refreshToken, err := generateJwtToken(application, user, nonce)
	if err != nil {
		panic(err)
	}

	if challenge == "null" {
		challenge = ""
	}

	token := &Token{
		Owner:         application.Owner,
		Name:          utils.GenerateId(),
		CreatedTime:   utils.GetCurrentTime(),
		Application:   application.Name,
		Organization:  user.Owner,
		User:          user.Name,
		Code:          utils.GenerateClientId(),
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		ExpiresIn:     application.ExpireInHours * 60,
		Scope:         scope,
		TokenType:     "Bearer",
		CodeChallenge: challenge,
	}
	AddToken(token)

	return &Code{
		Message: "",
		Code:    token.Code,
	}
}

func GetOAuthToken(grantType string, clientId string, clientSecret string, code string, verifier string) *TokenWrapper {
	application := GetApplicationByClientId(clientId)
	if application == nil {
		return &TokenWrapper{
			AccessToken: "error: invalid client_id",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}

	if grantType != "authorization_code" {
		return &TokenWrapper{
			AccessToken: "error: grant_type should be \"authorization_code\"",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}

	if code == "" {
		return &TokenWrapper{
			AccessToken: "error: code should not be empty",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}

	token := getTokenByCode(code)
	if token == nil {
		return &TokenWrapper{
			AccessToken: "error: invalid code",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}

	if application.Name != token.Application {
		return &TokenWrapper{
			AccessToken: "error: the token is for wrong application (client_id)",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}

	if application.ClientSecret != clientSecret {
		return &TokenWrapper{
			AccessToken: "error: invalid client_secret",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}
	if token.CodeChallenge != "" && pkceChallenge(verifier) != token.CodeChallenge {
		return &TokenWrapper{
			AccessToken: "error: incorrect code_verifier",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}

	tokenWrapper := &TokenWrapper{
		AccessToken:  token.AccessToken,
		IdToken:      token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		Scope:        token.Scope,
	}

	return tokenWrapper
}

func RefreshToken(grantType string, refreshToken string, scope string, clientId string, clientSecret string) *TokenWrapper {
	// check parameters
	if grantType != "refresh_token" {
		return &TokenWrapper{
			AccessToken: "error: grant_type should be \"refresh_token\"",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}
	application := GetApplicationByClientId(clientId)
	if application == nil {
		return &TokenWrapper{
			AccessToken: "error: invalid client_id",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}
	if application.ClientSecret != clientSecret {
		return &TokenWrapper{
			AccessToken: "error: invalid client_secret",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}
	// check whether the refresh token is valid, and has not expired.
	token := Token{RefreshToken: refreshToken}
	existed, err := adapter.Engine.Get(&token)
	if err != nil || !existed {
		return &TokenWrapper{
			AccessToken: "error: invalid refresh_token",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}

	cert := getCertByApplication(application)
	_, err = ParseJwtToken(refreshToken, cert)
	if err != nil {
		return &TokenWrapper{
			AccessToken: fmt.Sprintf("error: %s", err.Error()),
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}
	// generate a new token
	user := getUser(application.Organization, token.User)
	if user.IsForbidden {
		return &TokenWrapper{
			AccessToken: "error: the user is forbidden to sign in, please contact the administrator",
			TokenType:   "",
			ExpiresIn:   0,
			Scope:       "",
		}
	}
	newAccessToken, newRefreshToken, err := generateJwtToken(application, user, "")
	if err != nil {
		panic(err)
	}

	newToken := &Token{
		Owner:        application.Owner,
		Name:         utils.GenerateId(),
		CreatedTime:  utils.GetCurrentTime(),
		Application:  application.Name,
		Organization: user.Owner,
		User:         user.Name,
		Code:         utils.GenerateClientId(),
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    application.ExpireInHours * 60,
		Scope:        scope,
		TokenType:    "Bearer",
	}
	AddToken(newToken)

	tokenWrapper := &TokenWrapper{
		AccessToken:  token.AccessToken,
		IdToken:      token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		ExpiresIn:    token.ExpiresIn,
		Scope:        token.Scope,
	}

	return tokenWrapper
}

// PkceChallenge: base64-URL-encoded SHA256 hash of verifier, per rfc 7636
func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sum[:])
	return challenge
}
