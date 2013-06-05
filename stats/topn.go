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
	"fmt"
	"reflect"

//	"log"
)

// a datastructure to track the top N keys by value
type TopNContainer struct {
	size   int
	keys   []interface{}
	values []float64
}

func NewTopNContainer(size int) *TopNContainer {
	return &TopNContainer{
		size:   size,
		keys:   make([]interface{}, 0, size),
		values: make([]float64, 0, size),
	}
}

func (this *TopNContainer) Consider(key interface{}, invalue float64) {
	for i, val := range this.values {
		if invalue > val {
			this.Insert(key, invalue, i)
			return
		}
	}
	// if got to end, but not at capacity, add it to end
	if len(this.values) < this.size {
		this.Insert(key, invalue, len(this.values))
	}
}

func (this *TopNContainer) Insert(key interface{}, invalue float64, position int) {
	newActualSize := this.size
	if len(this.values) < newActualSize {
		newActualSize = len(this.values) + 1
	}

	newkeys := make([]interface{}, newActualSize, this.size)
	newvalues := make([]float64, newActualSize, this.size)
	for i := 0; i < newActualSize; i++ {
		if i < position {
			newkeys[i] = this.keys[i]
			newvalues[i] = this.values[i]
		} else if i == position {
			newkeys[i] = key
			newvalues[i] = invalue
		} else if i > position && i < this.size {
			newkeys[i] = this.keys[i-1]
			newvalues[i] = this.values[i-1]
		}
	}
	this.keys = newkeys
	this.values = newvalues
}

// returns 0 if the key isn't in the top N
func (this *TopNContainer) NumItemsWithKey(inkey interface{}) float64 {
	for i, key := range this.keys {
		if reflect.DeepEqual(key, inkey) {
			return this.values[i]
		}
	}
	return 0
}

func (this *TopNContainer) String() string {
	rv := ""
	for i := 0; i < len(this.values); i++ {
		if i != 0 {
			rv = rv + ", "
		}
		rv = rv + fmt.Sprintf("%v:%v", this.keys[i], this.values[i])
	}
	return rv
}
