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
	"time"
)

func (syncer *Syncer) syncUsers() {
	fmt.Printf("Running syncUsers()..\n")

	users, userMap := syncer.getUserMap()
	oUsers, oUserMap, err := syncer.getOriginalUserMap()
	if err != nil {
		fmt.Printf(err.Error())

		timestamp := time.Now().Format("2006-01-02 15:04:05")
		line := fmt.Sprintf("[%s] %s\n", timestamp, err.Error())
		updateSyncerErrorText(syncer, line)
		return
	}

	fmt.Printf("Users: %d, oUsers: %d\n", len(users), len(oUsers))

	var affiliationMap map[int]string
	if syncer.AffiliationTable != "" {
		_, affiliationMap = syncer.getAffiliationMap()
	}

	newUsers := []*User{}
	for _, oUser := range oUsers {
		id := oUser.Id
		if _, ok := userMap[id]; !ok {
			newUser := syncer.createUserFromOriginalUser(oUser, affiliationMap)
			fmt.Printf("New user: %v\n", newUser)
			newUsers = append(newUsers, newUser)
		} else {
			user := userMap[id]
			oHash := syncer.calculateHash(oUser)

			if user.Hash == user.PreHash {
				if user.Hash != oHash {
					updatedUser := syncer.createUserFromOriginalUser(oUser, affiliationMap)
					updatedUser.Hash = oHash
					updatedUser.PreHash = oHash
					syncer.updateUserForOriginalFields(updatedUser)
					fmt.Printf("Update from oUser to user: %v\n", updatedUser)
				}
			} else {
				if user.PreHash == oHash {
					updatedOUser := syncer.createOriginalUserFromUser(user)
					syncer.updateUser(updatedOUser)
					fmt.Printf("Update from user to oUser: %v\n", updatedOUser)

					// update preHash
					user.PreHash = user.Hash
					SetUserField(user, "pre_hash", user.PreHash)
				} else {
					if user.Hash == oHash {
						// update preHash
						user.PreHash = user.Hash
						SetUserField(user, "pre_hash", user.PreHash)
					} else {
						updatedUser := syncer.createUserFromOriginalUser(oUser, affiliationMap)
						updatedUser.Hash = oHash
						updatedUser.PreHash = oHash
						syncer.updateUserForOriginalFields(updatedUser)
						fmt.Printf("Update from oUser to user (2nd condition): %v\n", updatedUser)
					}
				}
			}
		}
	}
	AddUsersInBatch(newUsers)

	for _, user := range users {
		id := user.Id
		if _, ok := oUserMap[id]; !ok {
			newOUser := syncer.createOriginalUserFromUser(user)
			syncer.addUser(newOUser)
			fmt.Printf("New oUser: %v\n", newOUser)
		}
	}
}
