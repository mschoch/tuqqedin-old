//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package main

import (
	"reflect"
	"testing"

	"github.com/couchbaselabs/tuqqedin/plan"
)

func TestOptimizer(t *testing.T) {

	cheapPlan := plan.MockOperator{100}
	expensivePlan := plan.MockOperator{1000}

	optimizer := NewCouchbaseOptimizer()
	chosenPlan := optimizer.ChooseOptimalPlan([]plan.Operator{&cheapPlan, &expensivePlan})

	if !reflect.DeepEqual(chosenPlan, &cheapPlan) {
		t.Errorf("Expected plan %v, got %v", &cheapPlan, chosenPlan)
	}

}
