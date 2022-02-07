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
	"runtime"

	"github.com/bhojpur/dbm/pkg/core"
	"github.com/bhojpur/dbm/pkg/orm"

	"github.com/bhojpur/iam/pkg/conf"
	"github.com/bhojpur/iam/pkg/utils"
	websvr "github.com/bhojpur/web/pkg/engine"

	//_ "github.com/denisenkom/go-mssqldb" // db = mssql
	_ "github.com/go-sql-driver/mysql" // db = mysql
	//_ "github.com/lib/pq"                // db = postgres
)

var adapter *Adapter

func InitConfig() {
	err := websvr.LoadAppConfig("ini", "../conf/app.conf")
	if err != nil {
		panic(err)
	}

	InitAdapter(true)
}

func InitAdapter(createDatabase bool) {
	driverName, err := websvr.AppConfig.String("driverName")
	if err != nil {
		fmt.Errorf("driverName", err)
	}
	dbName, err := websvr.AppConfig.String("dbName")
	if err != nil {
		fmt.Errorf("dbName", err)
	}
	adapter = NewAdapter(driverName, conf.GetBhojpurConfDataSourceName(), dbName)
	if createDatabase {
		adapter.CreateDatabase()
	}
	adapter.createTable()
}

// Adapter represents the MySQL adapter for policy storage.
type Adapter struct {
	driverName     string
	dataSourceName string
	dbName         string
	Engine         *orm.Engine
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
	err := a.Engine.Close()
	if err != nil {
		panic(err)
	}
}

// NewAdapter is the constructor for Adapter.
func NewAdapter(driverName string, dataSourceName string, dbName string) *Adapter {
	a := &Adapter{}
	a.driverName = driverName
	a.dataSourceName = dataSourceName
	a.dbName = dbName

	// Open the DB, create it if not existed.
	a.open()

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a
}

func (a *Adapter) CreateDatabase() error {
	engine, err := orm.NewEngine(a.driverName, a.dataSourceName)
	if err != nil {
		return err
	}
	defer engine.Close()

	_, err = engine.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci", a.dbName))
	return err
}

func (a *Adapter) open() {
	dataSourceName := a.dataSourceName + a.dbName
	if a.driverName != "mysql" {
		dataSourceName = a.dataSourceName
	}

	engine, err := orm.NewEngine(a.driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	a.Engine = engine
}

func (a *Adapter) close() {
	_ = a.Engine.Close()
	a.Engine = nil
}

func (a *Adapter) createTable() {
	showSql, _ := websvr.AppConfig.Bool("showSql")
	a.Engine.ShowSQL(showSql)

	tableNamePrefix, err := websvr.AppConfig.String("tableNamePrefix")
	if err != nil {
		fmt.Errorf("tableNamePrefix", err)
	}
	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, tableNamePrefix)
	a.Engine.SetTableMapper(tbMapper)

	err = a.Engine.Sync2(new(Organization))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(User))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Role))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Permission))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Provider))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Application))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Resource))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Token))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(VerificationRecord))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Record))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Webhook))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Syncer))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Cert))
	if err != nil {
		panic(err)
	}

	err = a.Engine.Sync2(new(Ldap))
	if err != nil {
		panic(err)
	}
}

func GetSession(owner string, offset, limit int, field, value, sortField, sortOrder string) *orm.Session {
	session := adapter.Engine.Limit(limit, offset).Where("1=1")
	if owner != "" {
		session = session.And("owner=?", owner)
	}
	if field != "" && value != "" {
		session = session.And(fmt.Sprintf("%s like ?", utils.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
	}
	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}
	if sortOrder == "ascend" {
		session = session.Asc(utils.SnakeString(sortField))
	} else {
		session = session.Desc(utils.SnakeString(sortField))
	}
	return session
}
