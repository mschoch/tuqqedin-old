//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package stats

import (
	"math"
)

const NUM_FREQVALS = 10
const NUM_QUANTILES = 10

type PathStatistics struct {
	Rows               int
	DistinctValues     int
	MinValue           interface{}
	MaxValue           interface{}
	MostFrequentValues *TopNContainer
	Quantiles          []QuantileRange
}

func (this PathStatistics) NumFreqvals() int {
	return this.MostFrequentValues.size
}

func (this PathStatistics) NumQuantiles() int {
	return cap(this.Quantiles)
}

func DefaultPathStats(min, max interface{}) PathStatistics {
	return PathStatistics{
		Rows:               math.MaxInt32,
		DistinctValues:     math.MaxInt32,
		MinValue:           min,
		MaxValue:           max,
		MostFrequentValues: NewTopNContainer(NUM_FREQVALS),
		Quantiles:          make([]QuantileRange, 0, NUM_QUANTILES),
	}
}
