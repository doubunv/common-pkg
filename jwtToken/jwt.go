// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package jwtToken

import (
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

const (
	AccessTokenExpireDuration = 7 * 24 * time.Hour
	AdminToken                = 1
	UserToken                 = 2
)

var jwtSecret = []byte("kdd452-934sg4-l4d4q6")

type Claims struct {
	Uid          int64  `json:"uid"`
	Source       string `json:"source"`
	Sign         string `json:"sign"`
	BusinessCode string `json:"Business_code"`
	RoleId       int64  `json:"role_id"`
	TokenType    int    `json:"token_type"`
	jwt.StandardClaims
}

func GenerateToken(uid int64, source string, sign string, businessCode string, roleId int64, tokenType int) (accessToken string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(AccessTokenExpireDuration)
	claims := Claims{
		Uid:          uid,
		Source:       source,
		Sign:         sign,
		RoleId:       roleId,
		TokenType:    tokenType,
		BusinessCode: businessCode,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "test-user",
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil {
		return "", errors.Wrap(err, "failed to get accessToken")
	}

	return accessToken, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, errors.Wrap(err, "failed to parse token")
}

func ParseRefreshToken(aToken string) (claims *Claims, err error) {
	accessClaim, err := ParseToken(aToken)
	if err != nil {
		return claims, errors.WithMessage(err, "failed to parse accessToken")
	}

	if accessClaim.ExpiresAt > time.Now().Unix() {
		return accessClaim, nil
	}
	return accessClaim, errors.WithMessage(errors.New("Token time out"), "failed to parse refreshToken")
}
