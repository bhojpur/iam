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
	"time"

	"github.com/bhojpur/iam/pkg/utils"
	"github.com/bhopur/dbm/pkg/core"
)

type OriginalUser = User

func (syncer *Syncer) getOriginalUsers() ([]*OriginalUser, error) {
	sql := fmt.Sprintf("select * from %s", syncer.getTable())
	results, err := syncer.Adapter.Engine.QueryString(sql)
	if err != nil {
		return nil, err
	}

	return syncer.getOriginalUsersFromMap(results), nil
}

func (syncer *Syncer) getOriginalUserMap() ([]*OriginalUser, map[string]*OriginalUser, error) {
	users, err := syncer.getOriginalUsers()
	if err != nil {
		return users, nil, err
	}

	m := map[string]*OriginalUser{}
	for _, user := range users {
		m[user.Id] = user
	}
	return users, m, nil
}

func (syncer *Syncer) addUser(user *OriginalUser) (bool, error) {
	m := syncer.getMapFromOriginalUser(user)
	keyString, valueString := syncer.getSqlKeyValueStringFromMap(m)

	sql := fmt.Sprintf("insert into %s (%s) values (%s)", syncer.getTable(), keyString, valueString)
	res, err := syncer.Adapter.Engine.Exec(sql)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

/*func (syncer *Syncer) getOriginalColumns() []string {
	res := []string{}
	for _, tableColumn := range syncer.TableColumns {
		if tableColumn.BhojpurName != "Id" {
			res = append(res, tableColumn.Name)
		}
	}
	return res
}*/

func (syncer *Syncer) getBhojpurColumns() []string {
	res := []string{}
	for _, tableColumn := range syncer.TableColumns {
		if tableColumn.BhojpurName != "Id" {
			v := utils.CamelToSnakeCase(tableColumn.BhojpurName)
			res = append(res, v)
		}
	}
	return res
}

func (syncer *Syncer) updateUser(user *OriginalUser) (bool, error) {
	m := syncer.getMapFromOriginalUser(user)
	pkValue := m[syncer.TablePrimaryKey]
	delete(m, syncer.TablePrimaryKey)
	setString := syncer.getSqlSetStringFromMap(m)

	sql := fmt.Sprintf("update %s set %s where %s = %s", syncer.getTable(), setString, syncer.TablePrimaryKey, pkValue)
	res, err := syncer.Adapter.Engine.Exec(sql)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (syncer *Syncer) updateUserForOriginalFields(user *User) (bool, error) {
	owner, name := utils.GetOwnerAndNameFromId(user.GetId())
	oldUser := getUserById(owner, name)
	if oldUser == nil {
		return false, nil
	}

	if user.Avatar != oldUser.Avatar && user.Avatar != "" {
		user.PermanentAvatar = getPermanentAvatarUrl(user.Owner, user.Name, user.Avatar)
	}

	columns := syncer.getBhojpurColumns()
	columns = append(columns, "affiliation", "hash", "pre_hash")
	affected, err := adapter.Engine.ID(core.PK{oldUser.Owner, oldUser.Name}).Cols(columns...).Update(user)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (syncer *Syncer) calculateHash(user *OriginalUser) string {
	values := []string{}
	m := syncer.getMapFromOriginalUser(user)
	for _, tableColumn := range syncer.TableColumns {
		if tableColumn.IsHashed {
			values = append(values, m[tableColumn.Name])
		}
	}

	s := strings.Join(values, "|")
	return utils.GetMd5Hash(s)
}

func (syncer *Syncer) initAdapter() {
	if syncer.Adapter == nil {
		var dataSourceName string
		if syncer.DatabaseType == "mssql" {
			dataSourceName = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", syncer.User, syncer.Password, syncer.Host, syncer.Port, syncer.Database)
		} else {
			dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/", syncer.User, syncer.Password, syncer.Host, syncer.Port)
		}

		if !isCloudIntranet {
			dataSourceName = strings.ReplaceAll(dataSourceName, "dbi.", "db.")
		}

		syncer.Adapter = NewAdapter(syncer.DatabaseType, dataSourceName, syncer.Database)
	}
}

func RunSyncUsersJob() {
	syncers := GetSyncers("admin")
	for _, syncer := range syncers {
		addSyncerJob(syncer)
	}

	time.Sleep(time.Duration(1<<63 - 1))
}
