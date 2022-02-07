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

type Permission struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`
	DisplayName string `orm:"varchar(100)" json:"displayName"`

	Users []string `orm:"mediumtext" json:"users"`
	Roles []string `orm:"mediumtext" json:"roles"`

	ResourceType string   `orm:"varchar(100)" json:"resourceType"`
	Resources    []string `orm:"mediumtext" json:"resources"`
	Actions      []string `orm:"mediumtext" json:"actions"`
	Effect       string   `orm:"varchar(100)" json:"effect"`

	IsEnabled bool `json:"isEnabled"`
}

func GetPermissionCount(owner, field, value string) int {
	session := adapter.Engine.Where("owner=?", owner)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Permission{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetPermissions(owner string) []*Permission {
	permissions := []*Permission{}
	err := adapter.Engine.Desc("created_time").Find(&permissions, &Permission{Owner: owner})
	if err != nil {
		panic(err)
	}

	return permissions
}

func GetPaginationPermissions(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Permission {
	permissions := []*Permission{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&permissions)
	if err != nil {
		panic(err)
	}

	return permissions
}

func getPermission(owner string, name string) *Permission {
	if owner == "" || name == "" {
		return nil
	}

	permission := Permission{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&permission)
	if err != nil {
		panic(err)
	}

	if existed {
		return &permission
	} else {
		return nil
	}
}

func GetPermission(id string) *Permission {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getPermission(owner, name)
}

func UpdatePermission(id string, permission *Permission) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getPermission(owner, name) == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(permission)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddPermission(permission *Permission) bool {
	affected, err := adapter.Engine.Insert(permission)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeletePermission(permission *Permission) bool {
	affected, err := adapter.Engine.ID(core.PK{permission.Owner, permission.Name}).Delete(&Permission{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (permission *Permission) GetId() string {
	return fmt.Sprintf("%s/%s", permission.Owner, permission.Name)
}
