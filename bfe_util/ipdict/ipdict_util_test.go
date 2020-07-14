// Copyright (c) 2019 Baidu, Inc.
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

package ipdict

import (
	"net"
	"testing"
)

// normal case
func TestCheckIPPair_Case0(t *testing.T) {
	startIP := net.ParseIP("1.1.1.1")
	endIP := net.ParseIP("1.1.1.1")
	if err := checkIPPair(startIP, endIP); err != nil {
		t.Error(err.Error())
	}

	startIP = net.ParseIP("1::1")
	endIP = net.ParseIP("1::1")
	if err := checkIPPair(startIP, endIP); err != nil {
		t.Error(err.Error())
	}

}

// bad case
func TestCheckIPPair_Case1(t *testing.T) {
	startIP := net.ParseIP("1.1.1.1")
	endIP := net.ParseIP("1::1")
	if err := checkIPPair(startIP, endIP); err == nil {
		t.Error("TestCheckIPPair_Case0(): unexpected nil err")
	}

	startIP = net.ParseIP("0::1")
	endIP = net.ParseIP("1.1.1.1")
	if err := checkIPPair(startIP, endIP); err == nil {
		t.Error("TestCheckIPPair_Case0(): unexpected nil err")
	}

	startIP = net.ParseIP("1.1.1.1")[0:1]
	endIP = net.ParseIP("1.1.1.1")
	if err := checkIPPair(startIP, endIP); err == nil {
		t.Error("TestCheckIPPair_Case0(): unexpected nil err")
	}

	startIP = net.ParseIP("1.1.1.2")
	endIP = net.ParseIP("1.1.1.1")
	if err := checkIPPair(startIP, endIP); err == nil {
		t.Error("TestCheckIPPair_Case0(): unexpected nil err")
	}
}
