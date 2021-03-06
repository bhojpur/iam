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

type Provider struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`

	DisplayName   string `orm:"varchar(100)" json:"displayName"`
	Category      string `orm:"varchar(100)" json:"category"`
	Type          string `orm:"varchar(100)" json:"type"`
	SubType       string `orm:"varchar(100)" json:"subType"`
	Method        string `orm:"varchar(100)" json:"method"`
	ClientId      string `orm:"varchar(100)" json:"clientId"`
	ClientSecret  string `orm:"varchar(100)" json:"clientSecret"`
	ClientId2     string `orm:"varchar(100)" json:"clientId2"`
	ClientSecret2 string `orm:"varchar(100)" json:"clientSecret2"`

	Host    string `orm:"varchar(100)" json:"host"`
	Port    int    `json:"port"`
	Title   string `orm:"varchar(100)" json:"title"`
	Content string `orm:"varchar(1000)" json:"content"`

	RegionId     string `orm:"varchar(100)" json:"regionId"`
	SignName     string `orm:"varchar(100)" json:"signName"`
	TemplateCode string `orm:"varchar(100)" json:"templateCode"`
	AppId        string `orm:"varchar(100)" json:"appId"`

	Endpoint         string `orm:"varchar(1000)" json:"endpoint"`
	IntranetEndpoint string `orm:"varchar(100)" json:"intranetEndpoint"`
	Domain           string `orm:"varchar(100)" json:"domain"`
	Bucket           string `orm:"varchar(100)" json:"bucket"`

	Metadata               string `orm:"mediumtext" json:"metadata"`
	IdP                    string `orm:"mediumtext" json:"idP"`
	IssuerUrl              string `orm:"varchar(100)" json:"issuerUrl"`
	EnableSignAuthnRequest bool   `json:"enableSignAuthnRequest"`

	ProviderUrl string `orm:"varchar(200)" json:"providerUrl"`
}

func GetMaskedProvider(provider *Provider) *Provider {
	if provider == nil {
		return nil
	}

	if provider.ClientSecret != "" {
		provider.ClientSecret = "***"
	}
	if provider.ClientSecret2 != "" {
		provider.ClientSecret2 = "***"
	}

	return provider
}

func GetMaskedProviders(providers []*Provider) []*Provider {
	for _, provider := range providers {
		provider = GetMaskedProvider(provider)
	}
	return providers
}

func GetProviderCount(owner, field, value string) int {
	session := GetSession(owner, -1, -1, field, value, "", "")
	count, err := session.Count(&Provider{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetProviders(owner string) []*Provider {
	providers := []*Provider{}
	err := adapter.Engine.Desc("created_time").Find(&providers, &Provider{Owner: owner})
	if err != nil {
		panic(err)
	}

	return providers
}

func GetPaginationProviders(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Provider {
	providers := []*Provider{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&providers)
	if err != nil {
		panic(err)
	}

	return providers
}

func getProvider(owner string, name string) *Provider {
	if owner == "" || name == "" {
		return nil
	}

	provider := Provider{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&provider)
	if err != nil {
		panic(err)
	}

	if existed {
		return &provider
	} else {
		return nil
	}
}

func GetProvider(id string) *Provider {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getProvider(owner, name)
}

func GetDefaultHumanCheckProvider() *Provider {
	provider := Provider{Owner: "admin", Category: "HumanCheck"}
	existed, err := adapter.Engine.Get(&provider)
	if err != nil {
		panic(err)
	}

	if !existed {
		return nil
	}

	return &provider
}

func UpdateProvider(id string, provider *Provider) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getProvider(owner, name) == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(provider)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddProvider(provider *Provider) bool {
	affected, err := adapter.Engine.Insert(provider)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteProvider(provider *Provider) bool {
	affected, err := adapter.Engine.ID(core.PK{provider.Owner, provider.Name}).Delete(&Provider{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (p *Provider) GetId() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Name)
}
