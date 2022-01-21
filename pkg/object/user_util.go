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
	"reflect"
	"strings"

	"github.com/bhojpur/iam/pkg/idp"
	"xorm.io/core"
)

func GetUserByField(organizationName string, field string, value string) *User {
	if field == "" || value == "" {
		return nil
	}

	user := User{Owner: organizationName}
	existed, err := adapter.Engine.Where(fmt.Sprintf("%s=?", field), value).Get(&user)
	if err != nil {
		panic(err)
	}

	if existed {
		return &user
	} else {
		return nil
	}
}

func HasUserByField(organizationName string, field string, value string) bool {
	return GetUserByField(organizationName, field, value) != nil
}

func GetUserByFields(organization string, field string) *User {
	// check username
	user := GetUserByField(organization, "name", field)
	if user != nil {
		return user
	}

	// check email
	user = GetUserByField(organization, "email", field)
	if user != nil {
		return user
	}

	// check phone
	user = GetUserByField(organization, "phone", field)
	if user != nil {
		return user
	}

	// check ID card
	user = GetUserByField(organization, "id_card", field)
	if user != nil {
		return user
	}

	return nil
}

func SetUserField(user *User, field string, value string) bool {
	if field == "password" {
		organization := GetOrganizationByUser(user)
		user.UpdateUserPassword(organization)
		value = user.Password
	}

	affected, err := adapter.Engine.Table(user).ID(core.PK{user.Owner, user.Name}).Update(map[string]interface{}{field: value})
	if err != nil {
		panic(err)
	}

	user = getUser(user.Owner, user.Name)
	user.UpdateUserHash()
	_, err = adapter.Engine.ID(core.PK{user.Owner, user.Name}).Cols("hash").Update(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetUserField(user *User, field string) string {
	u := reflect.ValueOf(user)
	f := reflect.Indirect(u).FieldByName(field)
	return f.String()
}

func setUserProperty(user *User, field string, value string) {
	if value == "" {
		delete(user.Properties, field)
	} else {
		user.Properties[field] = value
	}
}

func SetUserOAuthProperties(organization *Organization, user *User, providerType string, userInfo *idp.UserInfo) bool {
	if userInfo.Id != "" {
		propertyName := fmt.Sprintf("oauth_%s_id", providerType)
		setUserProperty(user, propertyName, userInfo.Id)
	}
	if userInfo.Username != "" {
		propertyName := fmt.Sprintf("oauth_%s_username", providerType)
		setUserProperty(user, propertyName, userInfo.Username)
	}
	if userInfo.DisplayName != "" {
		propertyName := fmt.Sprintf("oauth_%s_displayName", providerType)
		setUserProperty(user, propertyName, userInfo.DisplayName)
		if user.DisplayName == "" {
			user.DisplayName = userInfo.DisplayName
		}
	}
	if userInfo.Email != "" {
		propertyName := fmt.Sprintf("oauth_%s_email", providerType)
		setUserProperty(user, propertyName, userInfo.Email)
		if user.Email == "" {
			user.Email = userInfo.Email
		}
	}
	if userInfo.AvatarUrl != "" {
		propertyName := fmt.Sprintf("oauth_%s_avatarUrl", providerType)
		setUserProperty(user, propertyName, userInfo.AvatarUrl)
		if user.Avatar == "" || user.Avatar == organization.DefaultAvatar {
			user.Avatar = userInfo.AvatarUrl
		}
	}

	affected := UpdateUserForAllFields(user.GetId(), user)
	return affected
}

func ClearUserOAuthProperties(user *User, providerType string) bool {
	for k := range user.Properties {
		prefix := fmt.Sprintf("oauth_%s_", providerType)
		if strings.HasPrefix(k, prefix) {
			delete(user.Properties, k)
		}
	}

	affected, err := adapter.Engine.ID(core.PK{user.Owner, user.Name}).Cols("properties").Update(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}
