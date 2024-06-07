// Copyright 2020 guylewin, guy@lewin.co.il
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"fmt"
	v1 "github.com/sunyihoo/frp/pkg/config/v1"
	"github.com/sunyihoo/frp/pkg/msg"
	"github.com/sunyihoo/frp/pkg/util/util"
	"slices"
)

type TokenAuthSetterVerifier struct {
	additionalAuthScopes []v1.AuthScope
	token                string
}

func NewTokenAuth(additionalAuthScopes []v1.AuthScope, token string) *TokenAuthSetterVerifier {
	return &TokenAuthSetterVerifier{
		additionalAuthScopes: additionalAuthScopes,
		token:                token,
	}
}

func (auth *TokenAuthSetterVerifier) VerifyLogin(m *msg.Login) error {
	if !util.ConstantTimeEqString(util.GetAuthKey(auth.token, m.Timestamp), m.PrivilegeKey) {
		return fmt.Errorf("token in login doesn't match token from configuration")
	}
	return nil
}

func (auth *TokenAuthSetterVerifier) VerifyPing(m *msg.Ping) error {
	if !slices.Contains(auth.additionalAuthScopes, v1.AuthScopeHeartBeats) {
		return nil
	}

	if !util.ConstantTimeEqString(util.GetAuthKey(auth.token, m.Timestamp), m.PrivilegeKey) {
		return fmt.Errorf("token in heartbeat doesn't match token from configuration")
	}
	return nil
}

func (auth *TokenAuthSetterVerifier) VerifyNewWorkConn(m *msg.NewWorkConn) error {
	if !slices.Contains(auth.additionalAuthScopes, v1.AuthScopeNewWorkConns) {
		return nil
	}
	if !util.ConstantTimeEqString(util.GetAuthKey(auth.token, m.Timestamp), m.PrivilegeKey) {
		return fmt.Errorf("token in NewWorkConn doesn't match token from configuration")
	}
	return nil
}
