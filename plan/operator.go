//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

import (
	"encoding/json"
	"fmt"
)

type Output interface{}
type OutputChannel chan Output

const NETWORK_COST = 10000
const MEMORY_COST = 100
const CPU_COST = 1

type Operator interface {
	Source() Operator
	GetOutputChannel() OutputChannel
	Run()
	Explain() map[string]interface{}
	Cancel()
	Cost() float64
	EstimatedRows() int
	TotalCost() float64
}

func OperatorToString(operator Operator) string {
	self := operator.Explain()
	bytes, err := json.MarshalIndent(self, "    ", "    ")
	if err != nil {
		return fmt.Sprintf("Error Serializing JSON: %v", err)
	}
	return string(bytes)
}
