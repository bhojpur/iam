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

import "fmt"

func getDbSyncerForUser(user *User) *Syncer {
	syncers := GetSyncers("admin")
	for _, syncer := range syncers {
		if syncer.Organization == user.Owner && syncer.IsEnabled && syncer.Type == "Database" {
			return syncer
		}
	}
	return nil
}

func getEnabledSyncerForOrganization(organization string) *Syncer {
	syncers := GetSyncers("admin")
	for _, syncer := range syncers {
		if syncer.Organization == organization && syncer.IsEnabled {
			return syncer
		}
	}
	return nil
}

func AddUserToOriginalDatabase(user *User) {
	syncer := getEnabledSyncerForOrganization(user.Owner)
	if syncer == nil {
		return
	}

	updatedOUser := syncer.createOriginalUserFromUser(user)
	syncer.addUser(updatedOUser)
	fmt.Printf("Add from user to oUser: %v\n", updatedOUser)
}

func UpdateUserToOriginalDatabase(user *User) {
	syncer := getEnabledSyncerForOrganization(user.Owner)
	if syncer == nil {
		return
	}

	newUser := GetUser(user.GetId())

	updatedOUser := syncer.createOriginalUserFromUser(newUser)
	syncer.updateUser(updatedOUser)
	fmt.Printf("Update from user to oUser: %v\n", updatedOUser)
}
