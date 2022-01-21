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
	"sync"
	"time"

	logsvr "github.com/bhojpur/logger/pkg/engine"
)

type LdapAutoSynchronizer struct {
	sync.Mutex
	ldapIdToStopChan map[string]chan struct{}
}

var globalLdapAutoSynchronizer *LdapAutoSynchronizer

func InitLdapAutoSynchronizer() {
	globalLdapAutoSynchronizer = NewLdapAutoSynchronizer()
	globalLdapAutoSynchronizer.LdapAutoSynchronizerStartUpAll()
}

func NewLdapAutoSynchronizer() *LdapAutoSynchronizer {
	return &LdapAutoSynchronizer{
		ldapIdToStopChan: make(map[string]chan struct{}),
	}
}

func GetLdapAutoSynchronizer() *LdapAutoSynchronizer {
	return globalLdapAutoSynchronizer
}

//start autosync for specified ldap, old existing autosync goroutine will be ceased
func (l *LdapAutoSynchronizer) StartAutoSync(ldapId string) error {
	l.Lock()
	defer l.Unlock()

	ldap := GetLdap(ldapId)
	if ldap == nil {
		return fmt.Errorf("ldap %s doesn't exist", ldapId)
	}
	if res, ok := l.ldapIdToStopChan[ldapId]; ok {
		res <- struct{}{}
		delete(l.ldapIdToStopChan, ldapId)
	}

	stopChan := make(chan struct{})
	l.ldapIdToStopChan[ldapId] = stopChan
	logsvr.Info(fmt.Sprintf("autoSync started for %s", ldap.Id))
	go l.syncRoutine(ldap, stopChan)
	return nil
}

func (l *LdapAutoSynchronizer) StopAutoSync(ldapId string) {
	l.Lock()
	defer l.Unlock()
	if res, ok := l.ldapIdToStopChan[ldapId]; ok {
		res <- struct{}{}
		delete(l.ldapIdToStopChan, ldapId)
	}
}

//autosync goroutine
func (l *LdapAutoSynchronizer) syncRoutine(ldap *Ldap, stopChan chan struct{}) {
	ticker := time.NewTicker(time.Duration(ldap.AutoSync) * time.Minute)
	defer ticker.Stop()
	for {
		UpdateLdapSyncTime(ldap.Id)
		//fetch all users
		conn, err := GetLdapConn(ldap.Host, ldap.Port, ldap.Admin, ldap.Passwd)
		if err != nil {
			logsvr.Warning(fmt.Sprintf("autoSync failed for %s, error %s", ldap.Id, err))
			continue
		}

		users, err := conn.GetLdapUsers(ldap.BaseDn)
		if err != nil {
			logsvr.Warning(fmt.Sprintf("autoSync failed for %s, error %s", ldap.Id, err))
			continue
		}
		existed, failed := SyncLdapUsers(ldap.Owner, LdapUsersToLdapRespUsers(users))
		if len(*failed) != 0 {
			logsvr.Warning(fmt.Sprintf("ldap autosync,%d new users,but %d user failed during :", len(users)-len(*existed)-len(*failed), len(*failed)), *failed)
		} else {
			logsvr.Info(fmt.Sprintf("ldap autosync success, %d new users, %d existing users", len(users)-len(*existed), len(*existed)))
		}
		select {
		case <-stopChan:
			logsvr.Info(fmt.Sprintf("autoSync goroutine for %s stopped", ldap.Id))
			return
		case <-ticker.C:
		}
	}

}

//start all autosync goroutine for existing ldap servers in each organizations
func (l *LdapAutoSynchronizer) LdapAutoSynchronizerStartUpAll() {
	organizations := []*Organization{}
	err := adapter.Engine.Desc("created_time").Find(&organizations)
	if err != nil {
		logsvr.Info("failed to Star up LdapAutoSynchronizer; ")
	}
	for _, org := range organizations {
		for _, ldap := range GetLdaps(org.Name) {
			if ldap.AutoSync != 0 {
				l.StartAutoSync(ldap.Id)
			}
		}
	}
}
