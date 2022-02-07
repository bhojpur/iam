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

type Payment struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`
	DisplayName string `orm:"varchar(100)" json:"displayName"`

	Provider     string `orm:"varchar(100)" json:"provider"`
	Type         string `orm:"varchar(100)" json:"type"`
	Organization string `orm:"varchar(100)" json:"organization"`
	User         string `orm:"varchar(100)" json:"user"`
	Good         string `orm:"varchar(100)" json:"good"`
	Amount       string `orm:"varchar(100)" json:"amount"`
	Currency     string `orm:"varchar(100)" json:"currency"`

	State string `orm:"varchar(100)" json:"state"`
}

func GetPaymentCount(owner, field, value string) int {
	session := GetSession(owner, -1, -1, field, value, "", "")
	count, err := session.Count(&Payment{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetPayments(owner string) []*Payment {
	payments := []*Payment{}
	err := adapter.Engine.Desc("created_time").Find(&payments, &Payment{Owner: owner})
	if err != nil {
		panic(err)
	}

	return payments
}

func GetPaginationPayments(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Payment {
	payments := []*Payment{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&payments)
	if err != nil {
		panic(err)
	}

	return payments
}

func getPayment(owner string, name string) *Payment {
	if owner == "" || name == "" {
		return nil
	}

	payment := Payment{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&payment)
	if err != nil {
		panic(err)
	}

	if existed {
		return &payment
	} else {
		return nil
	}
}

func GetPayment(id string) *Payment {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getPayment(owner, name)
}

func UpdatePayment(id string, payment *Payment) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getPayment(owner, name) == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(payment)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddPayment(payment *Payment) bool {
	affected, err := adapter.Engine.Insert(payment)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeletePayment(payment *Payment) bool {
	affected, err := adapter.Engine.ID(core.PK{payment.Owner, payment.Name}).Delete(&Payment{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (payment *Payment) GetId() string {
	return fmt.Sprintf("%s/%s", payment.Owner, payment.Name)
}
