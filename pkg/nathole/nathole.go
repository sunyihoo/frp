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

package nathole

import "net"

var (
	// SupportedModes
	// mode 0: simple detect mode, usually for both EasyNAT or HardNAT(Public Network)
	// a. receiver sends detect message with low TTL
	// b. sender sends normal detect message to receiver
	// c. receiver receives detect message and sends back a message to sender
	//
	// mod 1: For HardNAT & EasyNAT, send detect messages to multiple guessed ports.
	// Usually applicable to scenarios where port changes are regular.
	// Most of the steps are the same as mode 0, but EasyNAT is fixed as the receiver and will send detect messages
	// with low TTL to multiple guessed ports of the sender.
	//
	// mode 2: For HardNAT & EasyNAT, ports changes are not regular.
	// a. HardNAT machine will listen on multiple ports and send detect messages with low TTL to EasyNAT machine
	// b. EasyNAT machine will send detect messages to random ports of HardNAT machine.
	//
	// mode 3: For HardNAT & HardNAT, both changes in the ports regular.
	// Most of the steps are the same as mode 1, but the sender also needs to send detect messages to multiple guessed.
	// ports of the receiver.
	//
	// mode 4: For HardNAT & HardNAT, one of the changes in the ports is regular.
	// Regular port changes are usually on the sender side.
	// a. Receiver listens on multiple ports and sends detect messages with low TTL to the sender's guessed range ports.
	// b. Sender sends detect message to random ports of the receiver.
	SupportedModes = []int{DetectMode0, DetectMode1, DetectMode2, DetectMode3, DetectMode4}
	SupportedRoles = []string{DetectRoleSender, DetectRoleReceiver}

	DetectMode0        = 0
	DetectMode1        = 1
	DetectMode2        = 2
	DetectMode3        = 3
	DetectMode4        = 4
	DetectRoleSender   = "sender"
	DetectRoleReceiver = "receiver"
)

func parseIPs(addrs []string) []string {
	var ips []string
	for _, addr := range addrs {
		if ip, _, err := net.SplitHostPort(addr); err == nil {
			ips = append(ips, ip)
		}
	}
	return ips
}
