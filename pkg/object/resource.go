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

	"github.com/bhojpur/iam/pkg/utils"
	"github.com/bhopur/dbm/pkg/core"
)

type Resource struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`

	User        string `orm:"varchar(100)" json:"user"`
	Provider    string `orm:"varchar(100)" json:"provider"`
	Application string `orm:"varchar(100)" json:"application"`
	Tag         string `orm:"varchar(100)" json:"tag"`
	Parent      string `orm:"varchar(100)" json:"parent"`
	FileName    string `orm:"varchar(100)" json:"fileName"`
	FileType    string `orm:"varchar(100)" json:"fileType"`
	FileFormat  string `orm:"varchar(100)" json:"fileFormat"`
	FileSize    int    `json:"fileSize"`
	Url         string `orm:"varchar(1000)" json:"url"`
	Description string `orm:"varchar(1000)" json:"description"`
}

func GetResourceCount(owner, user, field, value string) int {
	session := adapter.Engine.Where("owner=? and user=?", owner, user)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Resource{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetResources(owner string, user string) []*Resource {
	if owner == "built-in" {
		owner = ""
		user = ""
	}

	resources := []*Resource{}
	err := adapter.Engine.Desc("created_time").Find(&resources, &Resource{Owner: owner, User: user})
	if err != nil {
		panic(err)
	}

	return resources
}

func GetPaginationResources(owner, user string, offset, limit int, field, value, sortField, sortOrder string) []*Resource {
	if owner == "built-in" {
		owner = ""
		user = ""
	}

	resources := []*Resource{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&resources, &Resource{User: user})
	if err != nil {
		panic(err)
	}

	return resources
}

func getResource(owner string, name string) *Resource {
	resource := Resource{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&resource)
	if err != nil {
		panic(err)
	}

	if existed {
		return &resource
	}

	return nil
}

func GetResource(id string) *Resource {
	owner, name := utils.GetOwnerAndNameFromIdNoCheck(id)
	return getResource(owner, name)
}

func UpdateResource(id string, resource *Resource) bool {
	owner, name := utils.GetOwnerAndNameFromIdNoCheck(id)
	if getResource(owner, name) == nil {
		return false
	}

	_, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(resource)
	if err != nil {
		panic(err)
	}

	//return affected != 0
	return true
}

func AddResource(resource *Resource) bool {
	affected, err := adapter.Engine.Insert(resource)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteResource(resource *Resource) bool {
	affected, err := adapter.Engine.ID(core.PK{resource.Owner, resource.Name}).Delete(&Resource{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (resource *Resource) GetId() string {
	return fmt.Sprintf("%s/%s", resource.Owner, resource.Name)
}

func AddOrUpdateResource(resource *Resource) bool {
	if getResource(resource.Owner, resource.Name) == nil {
		return AddResource(resource)
	} else {
		return UpdateResource(resource.GetId(), resource)
	}
}
