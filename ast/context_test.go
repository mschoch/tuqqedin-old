//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package ast

import (
	"reflect"
	"testing"
)

func TestNextPathElement(t *testing.T) {

	type Output struct {
		headString string
		headIndex  int
		restString string
		err        bool
	}

	tests := []struct {
		input  string
		output Output
	}{
		{"name", Output{"name", -1, "", false}},
		{"contact.firstName", Output{"contact", -1, "firstName", false}},
		{"contacts[0]", Output{"contacts", -1, "0]", false}},
		{"0]", Output{"", 0, "", false}},
		{"5].state", Output{"", 5, "state", false}},
		{"a]", Output{"", -1, "", true}},
	}

	for _, x := range tests {
		headString, headIndex, restString, err := NextPathElement(x.input)
		if err != nil && !x.output.err {
			t.Errorf("Expected no error, got error: %v", err)
		}
		if headString != x.output.headString {
			t.Errorf("Expected headString %v, got %v", x.output.headString, headString)
		}
		if headIndex != x.output.headIndex {
			t.Errorf("Expected headIndex %v, got %v", x.output.headIndex, headIndex)
		}
		if restString != x.output.restString {
			t.Errorf("Expected restString %v, got %v", x.output.restString, restString)
		}
	}

}

func TestContextGetPath(t *testing.T) {

	type Output struct {
		value interface{}
		err   bool
	}

	sampleContext := map[string]interface{}{
		"name": "will",
		"address": map[string]interface{}{
			"city": "New York",
		},
		"children": []interface{}{
			map[string]interface{}{
				"name": "bob",
			},
			map[string]interface{}{
				"name": "jane",
			},
		},
	}

	tests := []struct {
		input  string
		output Output
	}{
		{"name", Output{"will", false}},
		{"address.city", Output{"New York", false}},
		{"address.dne", Output{nil, false}},
		{"address.dne.dne", Output{nil, true}},
		{"children[0].name", Output{"bob", false}},
		{"children[1].name", Output{"jane", false}},
		{"children[1].name.xyz", Output{nil, true}},
		{"children[abc]", Output{nil, true}},
		{"address[0]", Output{nil, true}},
	}

	context := NewContext(sampleContext)

	for _, x := range tests {
		value, err := context.GetPath(x.input)
		if err != nil && !x.output.err {
			t.Errorf("Expected no error, got error: %v", err)
		}
		if !reflect.DeepEqual(value, x.output.value) {
			t.Errorf("Expected value %v, got %v", x.output.value, value)
		}

	}

}
