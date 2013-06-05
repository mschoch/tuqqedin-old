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
	"testing"
)

func TestTopNContainer(t *testing.T) {

	topTenList := NewTopNContainer(10)

	topTenList.Consider("marty", 100)

	if len(topTenList.keys) != 1 {
		t.Errorf("Expected length 1, got %v", len(topTenList.keys))
	}
	if len(topTenList.values) != 1 {
		t.Errorf("Expected length 1, got %v", len(topTenList.values))
	}

	topTenList.Consider("bob", 200)

	if len(topTenList.keys) != 2 {
		t.Errorf("Expected length 2, got %v", len(topTenList.keys))
	}
	if len(topTenList.values) != 2 {
		t.Errorf("Expected length 2, got %v", len(topTenList.values))
	}
	if topTenList.keys[0] != "bob" {
		t.Errorf("Expected bob on top, got %v", topTenList.keys[0])
	}
	if topTenList.keys[1] != "marty" {
		t.Errorf("Expected marty on bottom, got %v", topTenList.keys[1])
	}

	topTenList.Consider("third", 300)
	topTenList.Consider("fourth", 400)
	topTenList.Consider("fifth", 500)
	topTenList.Consider("sixth", 600)
	topTenList.Consider("sevent", 700)
	topTenList.Consider("eigth", 800)
	topTenList.Consider("ninth", 900)
	topTenList.Consider("tenth", 1000)
	topTenList.Consider("eleventh", 1100)

	if len(topTenList.keys) != 10 {
		t.Errorf("Expected length 2, got %v", len(topTenList.keys))
	}
	if len(topTenList.values) != 10 {
		t.Errorf("Expected length 2, got %v", len(topTenList.values))
	}
	if topTenList.keys[0] != "eleventh" {
		t.Errorf("Expected eleventh on top, got %v", topTenList.keys[0])
	}
	if topTenList.keys[9] != "bob" {
		t.Errorf("Expected bob on bottom, got %v", topTenList.keys[1])
	}

	t.Logf("Keys: %v", topTenList.keys)
	t.Logf("Values: %v", topTenList.values)
}
