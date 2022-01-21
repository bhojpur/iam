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

	"github.com/bhojpur/iam/pkg/cred"
	"github.com/bhojpur/iam/pkg/utils"
	"xorm.io/core"
)

type Organization struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	DisplayName        string `xorm:"varchar(100)" json:"displayName"`
	WebsiteUrl         string `xorm:"varchar(100)" json:"websiteUrl"`
	Favicon            string `xorm:"varchar(100)" json:"favicon"`
	PasswordType       string `xorm:"varchar(100)" json:"passwordType"`
	PasswordSalt       string `xorm:"varchar(100)" json:"passwordSalt"`
	PhonePrefix        string `xorm:"varchar(10)"  json:"phonePrefix"`
	DefaultAvatar      string `xorm:"varchar(100)" json:"defaultAvatar"`
	MasterPassword     string `xorm:"varchar(100)" json:"masterPassword"`
	EnableSoftDeletion bool   `json:"enableSoftDeletion"`
}

func GetOrganizationCount(owner, field, value string) int {
	session := adapter.Engine.Where("owner=?", owner)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Organization{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetOrganizations(owner string) []*Organization {
	organizations := []*Organization{}
	err := adapter.Engine.Desc("created_time").Find(&organizations, &Organization{Owner: owner})
	if err != nil {
		panic(err)
	}

	return organizations
}

func GetPaginationOrganizations(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Organization {
	organizations := []*Organization{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&organizations)
	if err != nil {
		panic(err)
	}

	return organizations
}

func getOrganization(owner string, name string) *Organization {
	if owner == "" || name == "" {
		return nil
	}

	organization := Organization{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&organization)
	if err != nil {
		panic(err)
	}

	if existed {
		return &organization
	}

	return nil
}

func GetOrganization(id string) *Organization {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getOrganization(owner, name)
}

func GetMaskedOrganization(organization *Organization) *Organization {
	if organization == nil {
		return nil
	}

	if organization.MasterPassword != "" {
		organization.MasterPassword = "***"
	}
	return organization
}

func GetMaskedOrganizations(organizations []*Organization) []*Organization {
	for _, organization := range organizations {
		organization = GetMaskedOrganization(organization)
	}
	return organizations
}

func UpdateOrganization(id string, organization *Organization) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getOrganization(owner, name) == nil {
		return false
	}

	if name == "built-in" {
		organization.Name = name
	}

	if name != organization.Name {
		applications := GetApplicationsByOrganizationName("admin", name)
		for _, application := range applications {
			application.Organization = organization.Name
			UpdateApplication(application.GetId(), application)
		}
	}

	if organization.MasterPassword != "" {
		credManager := cred.GetCredManager(organization.PasswordType)
		if credManager != nil {
			hashedPassword := credManager.GetHashedPassword(organization.MasterPassword, "", organization.PasswordSalt)
			organization.MasterPassword = hashedPassword
		}
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(organization)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddOrganization(organization *Organization) bool {
	affected, err := adapter.Engine.Insert(organization)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteOrganization(organization *Organization) bool {
	if organization.Name == "built-in" {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{organization.Owner, organization.Name}).Delete(&Organization{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetOrganizationByUser(user *User) *Organization {
	return getOrganization("admin", user.Owner)
}
