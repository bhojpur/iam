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
	"strings"

	"github.com/bhojpur/iam/pkg/utils"
	"github.com/bhopur/dbm/pkg/core"
)

type User struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`
	UpdatedTime string `orm:"varchar(100)" json:"updatedTime"`

	Id                string   `orm:"varchar(100) index" json:"id"`
	Type              string   `orm:"varchar(100)" json:"type"`
	Password          string   `orm:"varchar(100)" json:"password"`
	PasswordSalt      string   `orm:"varchar(100)" json:"passwordSalt"`
	DisplayName       string   `orm:"varchar(100)" json:"displayName"`
	Avatar            string   `orm:"varchar(500)" json:"avatar"`
	PermanentAvatar   string   `orm:"varchar(500)" json:"permanentAvatar"`
	Email             string   `orm:"varchar(100) index" json:"email"`
	Phone             string   `orm:"varchar(100) index" json:"phone"`
	Location          string   `orm:"varchar(100)" json:"location"`
	Address           []string `json:"address"`
	Affiliation       string   `orm:"varchar(100)" json:"affiliation"`
	Title             string   `orm:"varchar(100)" json:"title"`
	IdCardType        string   `orm:"varchar(100)" json:"idCardType"`
	IdCard            string   `orm:"varchar(100) index" json:"idCard"`
	Homepage          string   `orm:"varchar(100)" json:"homepage"`
	Bio               string   `orm:"varchar(100)" json:"bio"`
	Tag               string   `orm:"varchar(100)" json:"tag"`
	Region            string   `orm:"varchar(100)" json:"region"`
	Language          string   `orm:"varchar(100)" json:"language"`
	Gender            string   `orm:"varchar(100)" json:"gender"`
	Birthday          string   `orm:"varchar(100)" json:"birthday"`
	Education         string   `orm:"varchar(100)" json:"education"`
	Score             int      `json:"score"`
	Ranking           int      `json:"ranking"`
	IsDefaultAvatar   bool     `json:"isDefaultAvatar"`
	IsOnline          bool     `json:"isOnline"`
	IsAdmin           bool     `json:"isAdmin"`
	IsGlobalAdmin     bool     `json:"isGlobalAdmin"`
	IsForbidden       bool     `json:"isForbidden"`
	IsDeleted         bool     `json:"isDeleted"`
	SignupApplication string   `orm:"varchar(100)" json:"signupApplication"`
	Hash              string   `orm:"varchar(100)" json:"hash"`
	PreHash           string   `orm:"varchar(100)" json:"preHash"`

	CreatedIp      string `orm:"varchar(100)" json:"createdIp"`
	LastSigninTime string `orm:"varchar(100)" json:"lastSigninTime"`
	LastSigninIp   string `orm:"varchar(100)" json:"lastSigninIp"`

	Github   string `orm:"varchar(100)" json:"github"`
	Google   string `orm:"varchar(100)" json:"google"`
	QQ       string `orm:"qq varchar(100)" json:"qq"`
	WeChat   string `orm:"wechat varchar(100)" json:"wechat"`
	Facebook string `orm:"facebook varchar(100)" json:"facebook"`
	DingTalk string `orm:"dingtalk varchar(100)" json:"dingtalk"`
	Weibo    string `orm:"weibo varchar(100)" json:"weibo"`
	Gitee    string `orm:"gitee varchar(100)" json:"gitee"`
	LinkedIn string `orm:"linkedin varchar(100)" json:"linkedin"`
	Wecom    string `orm:"wecom varchar(100)" json:"wecom"`
	Lark     string `orm:"lark varchar(100)" json:"lark"`
	Gitlab   string `orm:"gitlab varchar(100)" json:"gitlab"`
	Apple    string `orm:"apple varchar(100)" json:"apple"`
	AzureAD  string `orm:"azuread varchar(100)" json:"azuread"`
	Slack    string `orm:"slack varchar(100)" json:"slack"`

	Ldap       string            `orm:"ldap varchar(100)" json:"ldap"`
	Properties map[string]string `json:"properties"`
}

func GetGlobalUserCount(field, value string) int {
	session := adapter.Engine.Where("1=1")
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&User{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetGlobalUsers() []*User {
	users := []*User{}
	err := adapter.Engine.Desc("created_time").Find(&users)
	if err != nil {
		panic(err)
	}

	return users
}

func GetPaginationGlobalUsers(offset, limit int, field, value, sortField, sortOrder string) []*User {
	users := []*User{}
	session := GetSession("", offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&users)
	if err != nil {
		panic(err)
	}

	return users
}

func GetUserCount(owner, field, value string) int {
	session := adapter.Engine.Where("owner=?", owner)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&User{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetOnlineUserCount(owner string, isOnline int) int {
	count, err := adapter.Engine.Where("is_online = ?", isOnline).Count(&User{Owner: owner})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetUsers(owner string) []*User {
	users := []*User{}
	err := adapter.Engine.Desc("created_time").Find(&users, &User{Owner: owner})
	if err != nil {
		panic(err)
	}

	return users
}

func GetSortedUsers(owner string, sorter string, limit int) []*User {
	users := []*User{}
	err := adapter.Engine.Desc(sorter).Limit(limit, 0).Find(&users, &User{Owner: owner})
	if err != nil {
		panic(err)
	}

	return users
}

func GetPaginationUsers(owner string, offset, limit int, field, value, sortField, sortOrder string) []*User {
	users := []*User{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&users)
	if err != nil {
		panic(err)
	}

	return users
}

func getUser(owner string, name string) *User {
	if owner == "" || name == "" {
		return nil
	}

	user := User{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&user)
	if err != nil {
		panic(err)
	}

	if existed {
		return &user
	} else {
		return nil
	}
}

func getUserById(owner string, id string) *User {
	if owner == "" || id == "" {
		return nil
	}

	user := User{Owner: owner, Id: id}
	existed, err := adapter.Engine.Get(&user)
	if err != nil {
		panic(err)
	}

	if existed {
		return &user
	} else {
		return nil
	}
}

func GetUserByEmail(owner string, email string) *User {
	if owner == "" || email == "" {
		return nil
	}

	user := User{Owner: owner, Email: email}
	existed, err := adapter.Engine.Get(&user)
	if err != nil {
		panic(err)
	}

	if existed {
		return &user
	} else {
		return nil
	}
}

func GetUser(id string) *User {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getUser(owner, name)
}

func GetUserNoCheck(id string) *User {
	owner, name := utils.GetOwnerAndNameFromIdNoCheck(id)
	return getUser(owner, name)
}

func GetMaskedUser(user *User) *User {
	if user == nil {
		return nil
	}

	if user.Password != "" {
		user.Password = "***"
	}
	return user
}

func GetMaskedUsers(users []*User) []*User {
	for _, user := range users {
		user = GetMaskedUser(user)
	}
	return users
}

func GetLastUser(owner string) *User {
	user := User{Owner: owner}
	existed, err := adapter.Engine.Desc("created_time", "id").Get(&user)
	if err != nil {
		panic(err)
	}

	if existed {
		return &user
	}

	return nil
}

func UpdateUser(id string, user *User, columns []string, isGlobalAdmin bool) bool {
	owner, name := utils.GetOwnerAndNameFromIdNoCheck(id)
	oldUser := getUser(owner, name)
	if oldUser == nil {
		return false
	}

	user.UpdateUserHash()

	if user.Avatar != oldUser.Avatar && user.Avatar != "" && user.PermanentAvatar != "*" {
		user.PermanentAvatar = getPermanentAvatarUrl(user.Owner, user.Name, user.Avatar)
	}

	if len(columns) == 0 {
		columns = []string{"owner", "display_name", "avatar",
			"location", "address", "region", "language", "affiliation", "title", "homepage", "bio", "score", "tag", "signup_application",
			"is_admin", "is_global_admin", "is_forbidden", "is_deleted", "hash", "is_default_avatar", "properties"}
	}
	if isGlobalAdmin {
		columns = append(columns, "name")
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).Cols(columns...).Update(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func UpdateUserForAllFields(id string, user *User) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	oldUser := getUser(owner, name)
	if oldUser == nil {
		return false
	}

	user.UpdateUserHash()

	if user.Avatar != oldUser.Avatar && user.Avatar != "" {
		user.PermanentAvatar = getPermanentAvatarUrl(user.Owner, user.Name, user.Avatar)
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddUser(user *User) bool {
	if user.Id == "" {
		user.Id = utils.GenerateId()
	}

	if user.Owner == "" || user.Name == "" {
		return false
	}

	organization := GetOrganizationByUser(user)
	user.UpdateUserPassword(organization)

	user.UpdateUserHash()
	user.PreHash = user.Hash

	user.PermanentAvatar = getPermanentAvatarUrl(user.Owner, user.Name, user.Avatar)

	affected, err := adapter.Engine.Insert(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddUsers(users []*User) bool {
	if len(users) == 0 {
		return false
	}

	//organization := GetOrganizationByUser(users[0])
	for _, user := range users {
		// this function is only used for syncer or batch upload, so no need to encrypt the password
		//user.UpdateUserPassword(organization)

		user.UpdateUserHash()
		user.PreHash = user.Hash

		user.PermanentAvatar = getPermanentAvatarUrl(user.Owner, user.Name, user.Avatar)
	}

	affected, err := adapter.Engine.Insert(users)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddUsersInBatch(users []*User) bool {
	batchSize := 1000

	if len(users) == 0 {
		return false
	}

	affected := false
	for i := 0; i < (len(users)-1)/batchSize+1; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(users) {
			end = len(users)
		}

		tmp := users[start:end]
		// TODO: save to log instead of standard output
		// fmt.Printf("Add users: [%d - %d].\n", start, end)
		if AddUsers(tmp) {
			affected = true
		}
	}

	return affected
}

func DeleteUser(user *User) bool {
	affected, err := adapter.Engine.ID(core.PK{user.Owner, user.Name}).Delete(&User{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func LinkUserAccount(user *User, field string, value string) bool {
	return SetUserField(user, field, value)
}

func (user *User) GetId() string {
	return fmt.Sprintf("%s/%s", user.Owner, user.Name)
}

func isUserIdGlobalAdmin(userId string) bool {
	return strings.HasPrefix(userId, "built-in/")
}
