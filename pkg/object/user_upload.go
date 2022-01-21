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
	"github.com/bhojpur/iam/pkg/utils"
	"github.com/bhojpur/iam/pkg/xlsx"
)

func getUserMap(owner string) map[string]*User {
	m := map[string]*User{}

	users := GetUsers(owner)
	for _, user := range users {
		m[user.GetId()] = user
	}

	return m
}

func parseLineItem(line *[]string, i int) string {
	if i >= len(*line) {
		return ""
	} else {
		return (*line)[i]
	}
}

func parseLineItemInt(line *[]string, i int) int {
	s := parseLineItem(line, i)
	return utils.ParseInt(s)
}

func parseLineItemBool(line *[]string, i int) bool {
	return parseLineItemInt(line, i) != 0
}

func UploadUsers(owner string, fileId string) bool {
	table := xlsx.ReadXlsxFile(fileId)

	oldUserMap := getUserMap(owner)
	newUsers := []*User{}
	for _, line := range table {
		if parseLineItem(&line, 0) == "" {
			continue
		}

		user := &User{
			Owner:             parseLineItem(&line, 0),
			Name:              parseLineItem(&line, 1),
			CreatedTime:       parseLineItem(&line, 2),
			UpdatedTime:       parseLineItem(&line, 3),
			Id:                parseLineItem(&line, 4),
			Type:              parseLineItem(&line, 5),
			Password:          parseLineItem(&line, 6),
			PasswordSalt:      parseLineItem(&line, 7),
			DisplayName:       parseLineItem(&line, 8),
			Avatar:            parseLineItem(&line, 9),
			PermanentAvatar:   "",
			Email:             parseLineItem(&line, 10),
			Phone:             parseLineItem(&line, 11),
			Location:          parseLineItem(&line, 12),
			Address:           []string{parseLineItem(&line, 13)},
			Affiliation:       parseLineItem(&line, 14),
			Title:             parseLineItem(&line, 15),
			IdCardType:        parseLineItem(&line, 16),
			IdCard:            parseLineItem(&line, 17),
			Homepage:          parseLineItem(&line, 18),
			Bio:               parseLineItem(&line, 19),
			Tag:               parseLineItem(&line, 20),
			Region:            parseLineItem(&line, 21),
			Language:          parseLineItem(&line, 22),
			Gender:            parseLineItem(&line, 23),
			Birthday:          parseLineItem(&line, 24),
			Education:         parseLineItem(&line, 25),
			Score:             parseLineItemInt(&line, 26),
			Ranking:           parseLineItemInt(&line, 27),
			IsDefaultAvatar:   false,
			IsOnline:          parseLineItemBool(&line, 28),
			IsAdmin:           parseLineItemBool(&line, 29),
			IsGlobalAdmin:     parseLineItemBool(&line, 30),
			IsForbidden:       parseLineItemBool(&line, 31),
			IsDeleted:         parseLineItemBool(&line, 32),
			SignupApplication: parseLineItem(&line, 33),
			Hash:              "",
			PreHash:           "",
			CreatedIp:         parseLineItem(&line, 34),
			LastSigninTime:    parseLineItem(&line, 35),
			LastSigninIp:      parseLineItem(&line, 36),
			Properties:        map[string]string{},
		}

		if _, ok := oldUserMap[user.GetId()]; !ok {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) == 0 {
		return false
	}
	return AddUsersInBatch(newUsers)
}
