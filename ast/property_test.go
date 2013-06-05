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

func TestEvaluateProperty(t *testing.T) {
	sampleDocument := map[string]interface{}{
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
		input  Expression
		output interface{}
	}{
		{NewProperty("name"), "will"},
	}

	context := NewContext(sampleDocument)

	for _, x := range tests {
		result, err := x.input.Evaluate(context)
		if err != nil {
			t.Fatalf("Error evaluating expression: %v", err)
		}
		if !reflect.DeepEqual(result, x.output) {
			t.Errorf("Expected %t %v, got %t %v", x.output, x.output, result, result)
		}
	}

}
