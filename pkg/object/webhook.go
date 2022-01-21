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
	"xorm.io/core"
)

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Webhook struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Organization string `xorm:"varchar(100) index" json:"organization"`

	Url            string    `xorm:"varchar(100)" json:"url"`
	Method         string    `xorm:"varchar(100)" json:"method"`
	ContentType    string    `xorm:"varchar(100)" json:"contentType"`
	Headers        []*Header `xorm:"mediumtext" json:"headers"`
	Events         []string  `xorm:"varchar(100)" json:"events"`
	IsUserExtended bool      `json:"isUserExtended"`
	IsEnabled      bool      `json:"isEnabled"`
}

func GetWebhookCount(owner, field, value string) int {
	session := adapter.Engine.Where("owner=?", owner)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Webhook{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetWebhooks(owner string) []*Webhook {
	webhooks := []*Webhook{}
	err := adapter.Engine.Desc("created_time").Find(&webhooks, &Webhook{Owner: owner})
	if err != nil {
		panic(err)
	}

	return webhooks
}

func GetPaginationWebhooks(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Webhook {
	webhooks := []*Webhook{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&webhooks)
	if err != nil {
		panic(err)
	}

	return webhooks
}

func getWebhooksByOrganization(organization string) []*Webhook {
	webhooks := []*Webhook{}
	err := adapter.Engine.Desc("created_time").Find(&webhooks, &Webhook{Organization: organization})
	if err != nil {
		panic(err)
	}

	return webhooks
}

func getWebhook(owner string, name string) *Webhook {
	if owner == "" || name == "" {
		return nil
	}

	webhook := Webhook{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&webhook)
	if err != nil {
		panic(err)
	}

	if existed {
		return &webhook
	} else {
		return nil
	}
}

func GetWebhook(id string) *Webhook {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getWebhook(owner, name)
}

func UpdateWebhook(id string, webhook *Webhook) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getWebhook(owner, name) == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(webhook)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddWebhook(webhook *Webhook) bool {
	affected, err := adapter.Engine.Insert(webhook)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteWebhook(webhook *Webhook) bool {
	affected, err := adapter.Engine.ID(core.PK{webhook.Owner, webhook.Name}).Delete(&Webhook{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (p *Webhook) GetId() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Name)
}
