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
	_ "embed"
	"fmt"
	"time"

	websvr "github.com/bhojpur/web/pkg/engine"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	*User
	Nonce string `json:"nonce,omitempty"`
	Tag   string `json:"tag,omitempty"`
	jwt.RegisteredClaims
}

type UserShort struct {
	Owner string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name  string `xorm:"varchar(100) notnull pk" json:"name"`
}

type ClaimsShort struct {
	*UserShort
	Nonce string `json:"nonce,omitempty"`
	jwt.RegisteredClaims
}

func getShortUser(user *User) *UserShort {
	res := &UserShort{
		Owner: user.Owner,
		Name:  user.Name,
	}
	return res
}

func getShortClaims(claims Claims) ClaimsShort {
	res := ClaimsShort{
		UserShort:        getShortUser(claims.User),
		Nonce:            claims.Nonce,
		RegisteredClaims: claims.RegisteredClaims,
	}
	return res
}

func generateJwtToken(application *Application, user *User, nonce string) (string, string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(application.ExpireInHours) * time.Hour)
	refreshExpireTime := nowTime.Add(time.Duration(application.RefreshExpireInHours) * time.Hour)

	user.Password = ""

	claims := Claims{
		User:  user,
		Nonce: nonce,
		// FIXME: A workaround for custom claim by reusing `tag` in user info
		Tag: user.Tag,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    websvr.AppConfig.String("origin"),
			Subject:   user.Id,
			Audience:  []string{application.ClientId},
			ExpiresAt: jwt.NewNumericDate(expireTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			ID:        "",
		},
	}

	var token *jwt.Token
	var refreshToken *jwt.Token

	// the JWT token length in "JWT-Empty" mode will be very short, as User object only has two properties: owner and name
	if application.TokenFormat == "JWT-Empty" {
		claimsShort := getShortClaims(claims)

		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claimsShort)
		claimsShort.ExpiresAt = jwt.NewNumericDate(refreshExpireTime)
		refreshToken = jwt.NewWithClaims(jwt.SigningMethodRS256, claimsShort)
	} else {
		token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		claims.ExpiresAt = jwt.NewNumericDate(refreshExpireTime)
		refreshToken = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	}

	cert := getCertByApplication(application)

	// RSA private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cert.PrivateKey))
	if err != nil {
		return "", "", err
	}

	token.Header["kid"] = cert.Name
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", "", err
	}
	refreshTokenString, err := refreshToken.SignedString(key)

	return tokenString, refreshTokenString, err
}

func ParseJwtToken(token string, cert *Cert) (*Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// RSA public key
		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert.PublicKey))
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	if t != nil {
		if claims, ok := t.Claims.(*Claims); ok && t.Valid {
			return claims, nil
		}
	}

	return nil, err
}
