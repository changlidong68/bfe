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

package mod_rewrite

import (
	"testing"
)

func TestReWriteTableSearch_1(t *testing.T) {
	config, err := ReWriteConfLoad("./testdata/rewrite_1.conf")
	if err != nil {
		t.Errorf("get err from ReWriteConfLoad():%s", err.Error())
		return
	}

	table := NewReWriteTable()
	table.Update(config)

	// search the table
	ruleList, ok := table.Search("pn")
	if !ok {
		t.Errorf("get err from Search():%s", err.Error())
		return
	}

	// check ruleList
	if len(*ruleList) != 1 {
		t.Error("len(ruleList) should be 1")
		return
	}
}
