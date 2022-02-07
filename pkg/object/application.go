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
	"fmt"

	"github.com/bhojpur/dbm/pkg/core"
	"github.com/bhojpur/iam/pkg/utils"
)

type Application struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`

	DisplayName         string          `orm:"varchar(100)" json:"displayName"`
	Logo                string          `orm:"varchar(100)" json:"logo"`
	HomepageUrl         string          `orm:"varchar(100)" json:"homepageUrl"`
	Description         string          `orm:"varchar(100)" json:"description"`
	Organization        string          `orm:"varchar(100)" json:"organization"`
	Cert                string          `orm:"varchar(100)" json:"cert"`
	EnablePassword      bool            `json:"enablePassword"`
	EnableSignUp        bool            `json:"enableSignUp"`
	EnableSigninSession bool            `json:"enableSigninSession"`
	EnableCodeSignin    bool            `json:"enableCodeSignin"`
	Providers           []*ProviderItem `orm:"mediumtext" json:"providers"`
	SignupItems         []*SignupItem   `orm:"varchar(1000)" json:"signupItems"`
	OrganizationObj     *Organization   `orm:"-" json:"organizationObj"`

	ClientId             string   `orm:"varchar(100)" json:"clientId"`
	ClientSecret         string   `orm:"varchar(100)" json:"clientSecret"`
	RedirectUris         []string `orm:"varchar(1000)" json:"redirectUris"`
	TokenFormat          string   `orm:"varchar(100)" json:"tokenFormat"`
	ExpireInHours        int      `json:"expireInHours"`
	RefreshExpireInHours int      `json:"refreshExpireInHours"`
	SignupUrl            string   `orm:"varchar(200)" json:"signupUrl"`
	SigninUrl            string   `orm:"varchar(200)" json:"signinUrl"`
	ForgetUrl            string   `orm:"varchar(200)" json:"forgetUrl"`
	AffiliationUrl       string   `orm:"varchar(100)" json:"affiliationUrl"`
	TermsOfUse           string   `orm:"varchar(100)" json:"termsOfUse"`
	SignupHtml           string   `orm:"mediumtext" json:"signupHtml"`
	SigninHtml           string   `orm:"mediumtext" json:"signinHtml"`
}

func GetApplicationCount(owner, field, value string) int {
	session := adapter.Engine.Where("owner=?", owner)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Application{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetApplications(owner string) []*Application {
	applications := []*Application{}
	err := adapter.Engine.Desc("created_time").Find(&applications, &Application{Owner: owner})
	if err != nil {
		panic(err)
	}

	return applications
}

func GetPaginationApplications(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Application {
	applications := []*Application{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&applications)
	if err != nil {
		panic(err)
	}

	return applications
}

func GetApplicationsByOrganizationName(owner string, organization string) []*Application {
	applications := []*Application{}
	err := adapter.Engine.Desc("created_time").Find(&applications, &Application{Owner: owner, Organization: organization})
	if err != nil {
		panic(err)
	}

	return applications
}

func getProviderMap(owner string) map[string]*Provider {
	providers := GetProviders(owner)
	m := map[string]*Provider{}
	for _, provider := range providers {
		//if provider.Category != "OAuth" {
		//	continue
		//}

		m[provider.Name] = GetMaskedProvider(provider)
	}
	return m
}

func extendApplicationWithProviders(application *Application) {
	m := getProviderMap(application.Owner)
	for _, providerItem := range application.Providers {
		if provider, ok := m[providerItem.Name]; ok {
			providerItem.Provider = provider
		}
	}
}

func extendApplicationWithOrg(application *Application) {
	organization := getOrganization(application.Owner, application.Organization)
	application.OrganizationObj = organization
}

func getApplication(owner string, name string) *Application {
	if owner == "" || name == "" {
		return nil
	}

	application := Application{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&application)
	if err != nil {
		panic(err)
	}

	if existed {
		extendApplicationWithProviders(&application)
		extendApplicationWithOrg(&application)
		return &application
	} else {
		return nil
	}
}

func GetApplicationByOrganizationName(organization string) *Application {
	application := Application{}
	existed, err := adapter.Engine.Where("organization=?", organization).Get(&application)
	if err != nil {
		panic(err)
	}

	if existed {
		extendApplicationWithProviders(&application)
		extendApplicationWithOrg(&application)
		return &application
	} else {
		return nil
	}
}

func GetApplicationByUser(user *User) *Application {
	if user.SignupApplication != "" {
		return getApplication("admin", user.SignupApplication)
	} else {
		return GetApplicationByOrganizationName(user.Owner)
	}
}

func GetApplicationByUserId(userId string) (*Application, *User) {
	var application *Application

	owner, name := utils.GetOwnerAndNameFromId(userId)
	if owner == "app" {
		application = getApplication("admin", name)
		return application, nil
	}

	user := GetUser(userId)
	application = GetApplicationByUser(user)

	return application, user
}

func GetApplicationByClientId(clientId string) *Application {
	application := Application{}
	existed, err := adapter.Engine.Where("client_id=?", clientId).Get(&application)
	if err != nil {
		panic(err)
	}

	if existed {
		extendApplicationWithProviders(&application)
		extendApplicationWithOrg(&application)
		return &application
	} else {
		return nil
	}
}

func GetApplication(id string) *Application {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getApplication(owner, name)
}

func GetMaskedApplication(application *Application, userId string) *Application {
	if isUserIdGlobalAdmin(userId) {
		return application
	}

	if application == nil {
		return nil
	}

	if application.ClientSecret != "" {
		application.ClientSecret = "***"
	}
	return application
}

func GetMaskedApplications(applications []*Application, userId string) []*Application {
	if isUserIdGlobalAdmin(userId) {
		return applications
	}

	for _, application := range applications {
		application = GetMaskedApplication(application, userId)
	}
	return applications
}

func UpdateApplication(id string, application *Application) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getApplication(owner, name) == nil {
		return false
	}

	if name == "app-built-in" {
		application.Name = name
	}

	for _, providerItem := range application.Providers {
		providerItem.Provider = nil
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(application)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddApplication(application *Application) bool {
	application.ClientId = utils.GenerateClientId()
	application.ClientSecret = utils.GenerateClientSecret()
	for _, providerItem := range application.Providers {
		providerItem.Provider = nil
	}

	affected, err := adapter.Engine.Insert(application)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteApplication(application *Application) bool {
	if application.Name == "app-built-in" {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{application.Owner, application.Name}).Delete(&Application{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (application *Application) GetId() string {
	return fmt.Sprintf("%s/%s", application.Owner, application.Name)
}
