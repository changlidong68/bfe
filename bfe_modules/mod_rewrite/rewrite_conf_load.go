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
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

import (
	"github.com/bfenetworks/bfe/bfe_basic/action"
	"github.com/bfenetworks/bfe/bfe_basic/condition"
)

type ReWriteRuleFile struct {
	Cond    *string         // condition for rewrite
	Actions []action.Action // list of actions
	Last    *bool           // if true, not to check the next rule in the list if
	// the condition is satisfied
}

type ReWriteRule struct {
	Cond    condition.Condition // condition for rewrite
	Actions []action.Action     // list of actions
	Last    bool                // if true, not to check the next rule in the list if
	// the condition is satisfied
}

type RuleFileList []ReWriteRuleFile
type RuleList []ReWriteRule

type ProductRulesFile map[string]*RuleFileList // product => list of rewrite rules
type ProductRules map[string]*RuleList         // product => list of rewrite rules

type ReWriteConfFile struct {
	Version *string // version of the config
	Config  *ProductRulesFile
}

type ReWriteConf struct {
	Version string       // version of the config
	Config  ProductRules // product rules for rewrite
}

func ReWriteRuleCheck(conf ReWriteRuleFile) error {
	// check Cond
	if conf.Cond == nil {
		return errors.New("no Cond")
	}

	// check Actions
	if conf.Actions == nil {
		return errors.New("no Actions")
	}

	// check Last
	if conf.Last == nil {
		return errors.New("no Last")
	}

	return nil
}

func RuleListCheck(conf *RuleFileList) error {
	for index, rule := range *conf {
		err := ReWriteRuleCheck(rule)
		if err != nil {
			return fmt.Errorf("ReWriteRule:%d, %s", index, err.Error())
		}
	}

	return nil
}

func ProductRulesCheck(conf *ProductRulesFile) error {
	for product, ruleList := range *conf {
		if ruleList == nil {
			return fmt.Errorf("no RuleList for product:%s", product)
		}

		err := RuleListCheck(ruleList)
		if err != nil {
			return fmt.Errorf("ProductRules:%s, %s", product, err.Error())
		}
	}

	return nil
}

func ReWriteConfCheck(conf ReWriteConfFile) error {
	var err error

	// check Version
	if conf.Version == nil {
		return errors.New("no Version")
	}

	// check Config
	if conf.Config == nil {
		return errors.New("no Config")
	}

	err = ProductRulesCheck(conf.Config)
	if err != nil {
		return fmt.Errorf("Config:%s", err.Error())
	}

	return nil
}

func ruleConvert(ruleFile ReWriteRuleFile) (ReWriteRule, error) {
	rule := ReWriteRule{}

	if ruleFile.Cond == nil {
		return rule, fmt.Errorf("cond not set")
	}
	cond, err := condition.Build(*ruleFile.Cond)
	if err != nil {
		return rule, err
	}
	rule.Cond = cond

	rule.Actions = ruleFile.Actions
	rule.Last = *ruleFile.Last
	return rule, nil
}

func ruleListConvert(ruleFileList *RuleFileList) (*RuleList, error) {
	ruleList := new(RuleList)
	*ruleList = make([]ReWriteRule, 0)

	for _, ruleFile := range *ruleFileList {
		rule, err := ruleConvert(ruleFile)
		if err != nil {
			return ruleList, err
		}
		*ruleList = append(*ruleList, rule)
	}

	return ruleList, nil
}

// ReWriteConfLoad loads config of rewrite from file.
func ReWriteConfLoad(filename string) (ReWriteConf, error) {
	var conf ReWriteConf
	var err error

	/* open the file    */
	file, err1 := os.Open(filename)

	if err1 != nil {
		return conf, err1
	}

	/* decode the file  */
	decoder := json.NewDecoder(file)

	var config ReWriteConfFile
	err = decoder.Decode(&config)
	file.Close()

	if err != nil {
		return conf, err
	}

	// check config
	err = ReWriteConfCheck(config)
	if err != nil {
		return conf, err
	}

	/* convert config   */
	conf.Version = *config.Version
	conf.Config = make(ProductRules)

	for product, ruleFileList := range *config.Config {
		ruleList, err := ruleListConvert(ruleFileList)
		if err != nil {
			return conf, err
		}
		conf.Config[product] = ruleList
	}

	return conf, nil
}
