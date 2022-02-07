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

type Cert struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`

	DisplayName     string `orm:"varchar(100)" json:"displayName"`
	Scope           string `orm:"varchar(100)" json:"scope"`
	Type            string `orm:"varchar(100)" json:"type"`
	CryptoAlgorithm string `orm:"varchar(100)" json:"cryptoAlgorithm"`
	BitSize         int    `json:"bitSize"`
	ExpireInYears   int    `json:"expireInYears"`

	PublicKey  string `orm:"mediumtext" json:"publicKey"`
	PrivateKey string `orm:"mediumtext" json:"privateKey"`
}

func GetMaskedCert(cert *Cert) *Cert {
	if cert == nil {
		return nil
	}

	return cert
}

func GetMaskedCerts(certs []*Cert) []*Cert {
	for _, cert := range certs {
		cert = GetMaskedCert(cert)
	}
	return certs
}

func GetCertCount(owner, field, value string) int {
	session := adapter.Engine.Where("owner=?", owner)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Cert{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetCerts(owner string) []*Cert {
	certs := []*Cert{}
	err := adapter.Engine.Desc("created_time").Find(&certs, &Cert{Owner: owner})
	if err != nil {
		panic(err)
	}

	return certs
}

func GetPaginationCerts(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Cert {
	certs := []*Cert{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&certs)
	if err != nil {
		panic(err)
	}

	return certs
}

func getCert(owner string, name string) *Cert {
	if owner == "" || name == "" {
		return nil
	}

	cert := Cert{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&cert)
	if err != nil {
		panic(err)
	}

	if existed {
		return &cert
	} else {
		return nil
	}
}

func GetCert(id string) *Cert {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getCert(owner, name)
}

func UpdateCert(id string, cert *Cert) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getCert(owner, name) == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(cert)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddCert(cert *Cert) bool {
	if cert.PublicKey == "" || cert.PrivateKey == "" {
		publicKey, privateKey := generateRsaKeys(cert.BitSize, cert.ExpireInYears, cert.Name, cert.Owner)
		cert.PublicKey = publicKey
		cert.PrivateKey = privateKey
	}

	affected, err := adapter.Engine.Insert(cert)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func DeleteCert(cert *Cert) bool {
	affected, err := adapter.Engine.ID(core.PK{cert.Owner, cert.Name}).Delete(&Cert{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func (p *Cert) GetId() string {
	return fmt.Sprintf("%s/%s", p.Owner, p.Name)
}

func getCertByApplication(application *Application) *Cert {
	if application.Cert != "" {
		return getCert("admin", application.Cert)
	} else {
		return GetDefaultCert()
	}
}

func GetDefaultCert() *Cert {
	return getCert("admin", "cert-built-in")
}
