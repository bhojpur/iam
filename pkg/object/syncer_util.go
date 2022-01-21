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
	"strconv"
	"strings"

	"github.com/bhojpur/iam/pkg/utils"
)

func (syncer *Syncer) getFullAvatarUrl(avatar string) string {
	if syncer.AvatarBaseUrl == "" {
		return avatar
	}

	if !strings.HasPrefix(avatar, "http") {
		return fmt.Sprintf("%s%s", syncer.AvatarBaseUrl, avatar)
	}
	return avatar
}

func (syncer *Syncer) getPartialAvatarUrl(avatar string) string {
	if strings.HasPrefix(avatar, syncer.AvatarBaseUrl) {
		return avatar[len(syncer.AvatarBaseUrl):]
	}
	return avatar
}

func (syncer *Syncer) createUserFromOriginalUser(originalUser *OriginalUser, affiliationMap map[int]string) *User {
	user := *originalUser
	user.Owner = syncer.Organization

	if user.Name == "" {
		user.Name = originalUser.Id
	}

	if user.CreatedTime == "" {
		user.CreatedTime = utils.GetCurrentTime()
	}

	if user.Type == "" {
		user.Type = "normal-user"
	}

	user.Avatar = syncer.getFullAvatarUrl(user.Avatar)

	if affiliationMap != nil {
		if originalUser.Score != 0 {
			affiliation, ok := affiliationMap[originalUser.Score]
			if !ok {
				panic(fmt.Sprintf("Affiliation not found: %d", originalUser.Score))
			}
			user.Affiliation = affiliation
		}
	}

	if user.Properties == nil {
		user.Properties = map[string]string{}
	}

	return &user
}

func (syncer *Syncer) createOriginalUserFromUser(user *User) *OriginalUser {
	originalUser := *user
	originalUser.Avatar = syncer.getPartialAvatarUrl(user.Avatar)
	return &originalUser
}

func (syncer *Syncer) setUserByKeyValue(user *User, key string, value string) {
	switch key {
	case "Name":
		user.Name = value
	case "CreatedTime":
		user.CreatedTime = value
	case "UpdatedTime":
		user.UpdatedTime = value
	case "Id":
		user.Id = value
	case "Type":
		user.Type = value
	case "Password":
		user.Password = value
	case "PasswordSalt":
		user.PasswordSalt = value
	case "DisplayName":
		user.DisplayName = value
	case "Avatar":
		user.Avatar = syncer.getPartialAvatarUrl(value)
	case "PermanentAvatar":
		user.PermanentAvatar = value
	case "Email":
		user.Email = value
	case "Phone":
		user.Phone = value
	case "Location":
		user.Location = value
	case "Address":
		user.Address = []string{value}
	case "Affiliation":
		user.Affiliation = value
	case "Title":
		user.Title = value
	case "IdCardType":
		user.IdCardType = value
	case "IdCard":
		user.IdCard = value
	case "Homepage":
		user.Homepage = value
	case "Bio":
		user.Bio = value
	case "Tag":
		user.Tag = value
	case "Region":
		user.Region = value
	case "Language":
		user.Language = value
	case "Gender":
		user.Gender = value
	case "Birthday":
		user.Birthday = value
	case "Education":
		user.Education = value
	case "Score":
		user.Score = utils.ParseInt(value)
	case "Ranking":
		user.Ranking = utils.ParseInt(value)
	case "IsDefaultAvatar":
		user.IsDefaultAvatar = utils.ParseBool(value)
	case "IsOnline":
		user.IsOnline = utils.ParseBool(value)
	case "IsAdmin":
		user.IsAdmin = utils.ParseBool(value)
	case "IsGlobalAdmin":
		user.IsGlobalAdmin = utils.ParseBool(value)
	case "IsForbidden":
		user.IsForbidden = utils.ParseBool(value)
	case "IsDeleted":
		user.IsDeleted = utils.ParseBool(value)
	case "CreatedIp":
		user.CreatedIp = value
	}
}

func (syncer *Syncer) getOriginalUsersFromMap(results []map[string]string) []*OriginalUser {
	users := []*OriginalUser{}
	for _, result := range results {
		originalUser := &OriginalUser{
			Address:    []string{},
			Properties: map[string]string{},
		}

		for _, tableColumn := range syncer.TableColumns {
			syncer.setUserByKeyValue(originalUser, tableColumn.BhojpurName, result[tableColumn.Name])
		}
		users = append(users, originalUser)
	}
	return users
}

func (syncer *Syncer) getMapFromOriginalUser(user *OriginalUser) map[string]string {
	m := map[string]string{}
	m["Name"] = user.Name
	m["CreatedTime"] = user.CreatedTime
	m["UpdatedTime"] = user.UpdatedTime
	m["Id"] = user.Id
	m["Type"] = user.Type
	m["Password"] = user.Password
	m["PasswordSalt"] = user.PasswordSalt
	m["DisplayName"] = user.DisplayName
	m["Avatar"] = syncer.getFullAvatarUrl(user.Avatar)
	m["PermanentAvatar"] = user.PermanentAvatar
	m["Email"] = user.Email
	m["Phone"] = user.Phone
	m["Location"] = user.Location
	m["Address"] = strings.Join(user.Address, "|")
	m["Affiliation"] = user.Affiliation
	m["Title"] = user.Title
	m["IdCardType"] = user.IdCardType
	m["IdCard"] = user.IdCard
	m["Homepage"] = user.Homepage
	m["Bio"] = user.Bio
	m["Tag"] = user.Tag
	m["Region"] = user.Region
	m["Language"] = user.Language
	m["Gender"] = user.Gender
	m["Birthday"] = user.Birthday
	m["Education"] = user.Education
	m["Score"] = strconv.Itoa(user.Score)
	m["Ranking"] = strconv.Itoa(user.Ranking)
	m["IsDefaultAvatar"] = utils.BoolToString(user.IsDefaultAvatar)
	m["IsOnline"] = utils.BoolToString(user.IsOnline)
	m["IsAdmin"] = utils.BoolToString(user.IsAdmin)
	m["IsGlobalAdmin"] = utils.BoolToString(user.IsGlobalAdmin)
	m["IsForbidden"] = utils.BoolToString(user.IsForbidden)
	m["IsDeleted"] = utils.BoolToString(user.IsDeleted)
	m["CreatedIp"] = user.CreatedIp

	m2 := map[string]string{}
	for _, tableColumn := range syncer.TableColumns {
		m2[tableColumn.Name] = m[tableColumn.BhojpurName]
	}

	return m2
}

func (syncer *Syncer) getSqlSetStringFromMap(m map[string]string) string {
	typeMap := syncer.getTableColumnsTypeMap()

	tokens := []string{}
	for k, v := range m {
		token := fmt.Sprintf("%s = %s", k, v)
		if typeMap[k] == "string" {
			token = fmt.Sprintf("%s = '%s'", k, v)
		}

		tokens = append(tokens, token)
	}
	return strings.Join(tokens, ", ")
}

func (syncer *Syncer) getSqlKeyValueStringFromMap(m map[string]string) (string, string) {
	typeMap := syncer.getTableColumnsTypeMap()

	keys := []string{}
	values := []string{}
	for k, v := range m {
		if typeMap[k] == "string" {
			v = fmt.Sprintf("'%s'", v)
		}

		keys = append(keys, k)
		values = append(values, v)
	}
	return strings.Join(keys, ", "), strings.Join(values, ", ")
}
