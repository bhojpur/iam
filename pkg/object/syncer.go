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

type TableColumn struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	BhojpurName string   `json:"bhojpurName"`
	IsHashed    bool     `json:"isHashed"`
	Values      []string `json:"values"`
}

type Syncer struct {
	Owner       string `orm:"varchar(100) notnull pk" json:"owner"`
	Name        string `orm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `orm:"varchar(100)" json:"createdTime"`

	Organization string `orm:"varchar(100)" json:"organization"`
	Type         string `orm:"varchar(100)" json:"type"`

	Host             string         `orm:"varchar(100)" json:"host"`
	Port             int            `json:"port"`
	User             string         `orm:"varchar(100)" json:"user"`
	Password         string         `orm:"varchar(100)" json:"password"`
	DatabaseType     string         `orm:"varchar(100)" json:"databaseType"`
	Database         string         `orm:"varchar(100)" json:"database"`
	Table            string         `orm:"varchar(100)" json:"table"`
	TablePrimaryKey  string         `orm:"varchar(100)" json:"tablePrimaryKey"`
	TableColumns     []*TableColumn `orm:"mediumtext" json:"tableColumns"`
	AffiliationTable string         `orm:"varchar(100)" json:"affiliationTable"`
	AvatarBaseUrl    string         `orm:"varchar(100)" json:"avatarBaseUrl"`
	ErrorText        string         `orm:"mediumtext" json:"errorText"`
	SyncInterval     int            `json:"syncInterval"`
	IsEnabled        bool           `json:"isEnabled"`

	Adapter *Adapter `orm:"-" json:"-"`
}

func GetSyncerCount(owner, field, value string) int {
	session := adapter.Engine.Where("owner=?", owner)
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	count, err := session.Count(&Syncer{})
	if err != nil {
		panic(err)
	}

	return int(count)
}

func GetSyncers(owner string) []*Syncer {
	syncers := []*Syncer{}
	err := adapter.Engine.Desc("created_time").Find(&syncers, &Syncer{Owner: owner})
	if err != nil {
		panic(err)
	}

	return syncers
}

func GetPaginationSyncers(owner string, offset, limit int, field, value, sortField, sortOrder string) []*Syncer {
	syncers := []*Syncer{}
	session := GetSession(owner, offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&syncers)
	if err != nil {
		panic(err)
	}

	return syncers
}

func getSyncer(owner string, name string) *Syncer {
	if owner == "" || name == "" {
		return nil
	}

	syncer := Syncer{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&syncer)
	if err != nil {
		panic(err)
	}

	if existed {
		return &syncer
	} else {
		return nil
	}
}

func GetSyncer(id string) *Syncer {
	owner, name := utils.GetOwnerAndNameFromId(id)
	return getSyncer(owner, name)
}

func GetMaskedSyncer(syncer *Syncer) *Syncer {
	if syncer == nil {
		return nil
	}

	if syncer.Password != "" {
		syncer.Password = "***"
	}
	return syncer
}

func GetMaskedSyncers(syncers []*Syncer) []*Syncer {
	for _, syncer := range syncers {
		syncer = GetMaskedSyncer(syncer)
	}
	return syncers
}

func UpdateSyncer(id string, syncer *Syncer) bool {
	owner, name := utils.GetOwnerAndNameFromId(id)
	if getSyncer(owner, name) == nil {
		return false
	}

	affected, err := adapter.Engine.ID(core.PK{owner, name}).AllCols().Update(syncer)
	if err != nil {
		panic(err)
	}

	if affected == 1 {
		addSyncerJob(syncer)
	}

	return affected != 0
}

func updateSyncerErrorText(syncer *Syncer, line string) bool {
	s := getSyncer(syncer.Owner, syncer.Name)
	if s == nil {
		return false
	}

	s.ErrorText = s.ErrorText + line

	affected, err := adapter.Engine.ID(core.PK{s.Owner, s.Name}).Cols("error_text").Update(s)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddSyncer(syncer *Syncer) bool {
	affected, err := adapter.Engine.Insert(syncer)
	if err != nil {
		panic(err)
	}

	if affected == 1 {
		addSyncerJob(syncer)
	}

	return affected != 0
}

func DeleteSyncer(syncer *Syncer) bool {
	affected, err := adapter.Engine.ID(core.PK{syncer.Owner, syncer.Name}).Delete(&Syncer{})
	if err != nil {
		panic(err)
	}

	if affected == 1 {
		deleteSyncerJob(syncer)
	}

	return affected != 0
}

func (syncer *Syncer) GetId() string {
	return fmt.Sprintf("%s/%s", syncer.Owner, syncer.Name)
}

func (syncer *Syncer) getTableColumnsTypeMap() map[string]string {
	m := map[string]string{}
	for _, tableColumn := range syncer.TableColumns {
		m[tableColumn.Name] = tableColumn.Type
	}
	return m
}

func (syncer *Syncer) getTable() string {
	if syncer.DatabaseType == "mssql" {
		return fmt.Sprintf("[%s]", syncer.Table)
	} else {
		return syncer.Table
	}
}
