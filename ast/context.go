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
	"fmt"
	//"log"
	"strconv"

	"strings"
)

type Context interface {
	GetPath(path string) (interface{}, error)
}

type RelativeContext struct {
	raw        map[string]interface{}
	relativeTo string
}

func NewContext(raw map[string]interface{}) *RelativeContext {
	return &RelativeContext{
		raw:        raw,
		relativeTo: "",
	}
}

func (this *RelativeContext) GetSubContext(path string) Context {
	this.relativeTo = this.relativeTo + "." + path
	return this
}

func (this *RelativeContext) GetPath(path string) (interface{}, error) {
	// TODO add support for absolute paths
	accessPath := path
	if this.relativeTo != "" {
		accessPath = this.relativeTo + "." + path
	}

	var currentPath = ""
	var curr interface{}
	curr = this.raw
	for accessPath != "" {
		headPath, headIndex, restPath, err := NextPathElement(accessPath)
		if err != nil {
			return nil, err
		}
		if headPath != "" {
			switch inside := curr.(type) {
			case map[string]interface{}:
				curr = inside[headPath]
				accessPath = restPath
				currentPath = currentPath + "." + headPath
			default:
				return nil, fmt.Errorf("Cannot access property %v within %v it is not an object it is %T", headPath, currentPath, inside)
			}
		} else if headIndex != -1 {
			switch inside := curr.(type) {
			case []interface{}:
				curr = inside[headIndex]
				accessPath = restPath
				currentPath = currentPath + "." + headPath
			default:
				return nil, fmt.Errorf("Cannot access index %v within %v it is not an array it is %T", headPath, currentPath, inside)
			}
		} else {
			return nil, fmt.Errorf("Unexpected state")
		}

	}
	return curr, nil
}

func NextPathElement(path string) (headPath string, headIndex int, rest string, err error) {
	dotIndex := strings.Index(path, ".")
	lbIndex := strings.Index(path, "[")
	rbIndex := strings.Index(path, "]")

	// first look for dangling rb
	if rbIndex > 0 && (dotIndex < 0 || rbIndex < dotIndex) && (lbIndex < 0 || rbIndex < lbIndex) {
		index := path[0:rbIndex]
		indexInt, err := strconv.Atoi(index)
		if err != nil {
			return "", -1, "", err
		}
		if len(path) > rbIndex+1 && string(path[rbIndex+1]) == "." {
			// skip over the dot too
			rbIndex = rbIndex + 1
		}
		return "", indexInt, path[rbIndex+1:], nil
	} else if dotIndex > 0 && (lbIndex < 0 || dotIndex < lbIndex) {
		return path[0:dotIndex], -1, path[dotIndex+1:], nil
	} else if lbIndex > 0 {
		return path[0:lbIndex], -1, path[lbIndex+1:], nil
	}

	return path, -1, "", nil
}
