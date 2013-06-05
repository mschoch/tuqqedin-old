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
	"log"

	"github.com/couchbaselabs/tuqqedin/plan"
)

type Optimizer interface {
	ChooseOptimalPlan([]plan.Operator) plan.Operator
}

type CouchbaseOptimizer struct {
}

func NewCouchbaseOptimizer() *CouchbaseOptimizer {
	return &CouchbaseOptimizer{}
}

func (this *CouchbaseOptimizer) ChooseOptimalPlan(plans []plan.Operator) plan.Operator {
	cheapest := plans[0]
	for i, plan := range plans {
		planCost := plan.TotalCost()
		log.Printf("Plan %d has cost: %v", i, planCost)
		if planCost < cheapest.TotalCost() {
			cheapest = plan
		}
	}

	return cheapest
}
