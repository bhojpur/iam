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
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/bhojpur/dbm/pkg/core"
	"github.com/bhojpur/iam/pkg/utils"
	websvr "github.com/bhojpur/web/pkg/engine"
)

type VerificationRecord struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`

	RemoteAddr string `orm:"varchar(100)"`
	Type       string `orm:"varchar(10)"`
	User       string `orm:"varchar(100) notnull"`
	Provider   string `orm:"varchar(100) notnull"`
	Receiver   string `orm:"varchar(100) notnull"`
	Code       string `orm:"varchar(10) notnull"`
	Time       int64  `orm:"notnull"`
	IsUsed     bool
}

func SendVerificationCodeToEmail(organization *Organization, user *User, provider *Provider, remoteAddr string, dest string) error {
	if provider == nil {
		return fmt.Errorf("Please set an Email provider first")
	}

	sender := organization.DisplayName
	title := provider.Title
	code := getRandomCode(5)
	// "You have requested a verification code at Bhojpur IAM. Here is your code: %s, please enter in 5 minutes."
	content := fmt.Sprintf(provider.Content, code)

	if err := AddToVerificationRecord(user, provider, remoteAddr, provider.Category, dest, code); err != nil {
		return err
	}

	return SendEmail(provider, title, content, dest, sender)
}

func SendVerificationCodeToPhone(organization *Organization, user *User, provider *Provider, remoteAddr string, dest string) error {
	if provider == nil {
		return errors.New("Please set a SMS provider first")
	}

	code := getRandomCode(5)
	if err := AddToVerificationRecord(user, provider, remoteAddr, provider.Category, dest, code); err != nil {
		return err
	}

	return SendSms(provider, code, dest)
}

func AddToVerificationRecord(user *User, provider *Provider, remoteAddr, recordType, dest, code string) error {
	var record VerificationRecord
	record.RemoteAddr = remoteAddr
	record.Type = recordType
	if user != nil {
		record.User = user.GetId()
	}
	has, err := adapter.Engine.Desc("created_time").Get(&record)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	if has && now-record.Time < 60 {
		return errors.New("You can only send one code in 60s.")
	}

	record.Owner = provider.Owner
	record.Name = utils.GenerateId()
	record.CreatedTime = utils.GetCurrentTime()
	if user != nil {
		record.User = user.GetId()
	}
	record.Provider = provider.Name

	record.Receiver = dest
	record.Code = code
	record.Time = now
	record.IsUsed = false

	_, err = adapter.Engine.Insert(record)
	if err != nil {
		return err
	}

	return nil
}

func getVerificationRecord(dest string) *VerificationRecord {
	var record VerificationRecord
	record.Receiver = dest
	has, err := adapter.Engine.Desc("time").Where("is_used = false").Get(&record)
	if err != nil {
		panic(err)
	}
	if !has {
		return nil
	}
	return &record
}

func CheckVerificationCode(dest, code string) string {
	record := getVerificationRecord(dest)

	if record == nil {
		return "Code has not been sent yet!"
	}

	timeout, err := websvr.AppConfig.Int64("verificationCodeTimeout")
	if err != nil {
		panic(err)
	}

	now := time.Now().Unix()
	if now-record.Time > timeout*60 {
		return fmt.Sprintf("You should verify your code in %d min!", timeout)
	}

	if record.Code != code {
		return "Wrong code!"
	}

	return ""
}

func DisableVerificationCode(dest string) {
	record := getVerificationRecord(dest)
	if record == nil {
		return
	}

	record.IsUsed = true
	_, err := adapter.Engine.ID(core.PK{record.Owner, record.Name}).AllCols().Update(record)
	if err != nil {
		panic(err)
	}
}

// from VnfNode/object/validateCode.go line 116
var stdNums = []byte("0123456789")

func getRandomCode(length int) string {
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, stdNums[r.Intn(len(stdNums))])
	}
	return string(result)
}
