// Copyright 2023 The frp Authors
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
	"github.com/sunyihoo/frp/pkg/msg"
)

var AlwaysPassVerifier = &alwaysPass{}

var _ Verifier = &alwaysPass{}

type alwaysPass struct{}

func (*alwaysPass) VerifyLogin(*msg.Login) error { return nil }

func (*alwaysPass) VerifyPing(*msg.Ping) error { return nil }

func (*alwaysPass) VerifyNewWorkConn(*msg.NewWorkConn) error { return nil }