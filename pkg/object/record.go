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
	ctxsvr "github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

var logPostOnly bool

func init() {
	var err error
	logPostOnly, err = websvr.AppConfig.Bool("logPostOnly")
	if err != nil {
		//panic(err)
	}
}

type Record struct {
	Id int `orm:"int notnull pk autoincr" json:"id"`

	Owner       string `orm:"varchar(100) index" json:"owner"`
	Name        string `orm:"varchar(100) index" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`

	Organization string `orm:"varchar(100)" json:"organization"`
	ClientIp     string `orm:"varchar(100)" json:"clientIp"`
	User         string `orm:"varchar(100)" json:"user"`
	Method       string `orm:"varchar(100)" json:"method"`
	RequestUri   string `orm:"varchar(1000)" json:"requestUri"`
	Action       string `orm:"varchar(1000)" json:"action"`

	ExtendedUser *User `orm:"-" json:"extendedUser"`

	IsTriggered bool `json:"isTriggered"`
}

func NewRecord(ctx *ctxsvr.Context) *Record {
	ip := strings.Replace(utils.GetIPFromRequest(ctx.Request), ": ", "", -1)
	action := strings.Replace(ctx.Request.URL.Path, "/api/", "", -1)
	requestUri := utils.FilterQuery(ctx.Request.RequestURI, []string{"accessToken"})
	if len(requestUri) > 1000 {
		requestUri = requestUri[0:1000]
	}

	record := Record{
		Name:        utils.GenerateId(),
		CreatedTime: utils.GetCurrentTime(),
		ClientIp:    ip,
		User:        "",
		Method:      ctx.Request.Method,
		RequestUri:  requestUri,
		Action:      action,
		IsTriggered: false,
	}
	return &record
}

func AddRecord(record *Record) bool {
	if logPostOnly {
		if record.Method == "GET" {
			return false
		}
	}

	if record.Organization == "app" {
		return false
	}

	record.Owner = record.Organization

	errWebhook := SendWebhooks(record)
	if errWebhook == nil {
		record.IsTriggered = true
	} else {
		fmt.Println(errWebhook)
	}

	affected, err := adapter.Engine.Insert(record)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetRecordCount(field, value string) int {
	session := adapter.Engine.Where("1=1")
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Record{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetRecords() []*Record {
	records := []*Record{}
	err := adapter.Engine.Desc("id").Find(&records)
	if err != nil {
		panic(err)
	}

	return records
}

func GetPaginationRecords(offset, limit int, field, value, sortField, sortOrder string) []*Record {
	records := []*Record{}
	session := GetSession("", offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&records)
	if err != nil {
		panic(err)
	}

	return records
}

func GetRecordsByField(record *Record) []*Record {
	records := []*Record{}
	err := adapter.Engine.Find(&records, record)
	if err != nil {
		panic(err)
	}

	return records
}

func SendWebhooks(record *Record) error {
	webhooks := getWebhooksByOrganization(record.Organization)
	for _, webhook := range webhooks {
		if !webhook.IsEnabled {
			continue
		}

		matched := false
		for _, event := range webhook.Events {
			if record.Action == event {
				matched = true
				break
			}
		}

		if matched {
			if webhook.IsUserExtended {
				user := getUser(record.Organization, record.User)
				record.ExtendedUser = user
			}

			err := sendWebhook(webhook, record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
