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
	//	"fmt"
	"log"
	"strings"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/couchbaselabs/tuqqedin/ast"
	"github.com/couchbaselabs/tuqqedin/datasource"
)

var MIN_KEY interface{}
var MAX_KEY = map[string]interface{}{} //FIXME better approximation of max

// a key in Couchbase is at most 250 bytes
// in a view, items with the same key are sorted by id (using basic memcmp comparison)
// so 251 bytes of 0xff is higher than any valid key in couchbase
var MAX_ID = couchbase.DocId(strings.Repeat(string([]byte{0xff}), 251))

// an empty string is an invalid id, but will sort before any other key
var MIN_ID = couchbase.DocId("")

var MIN_LOCATION = NewViewLocationGreatherThan(MIN_KEY, true)
var MAX_LOCATION = NewViewLocationLessThan(MAX_KEY, true)

type ViewLocation struct {
	Key   interface{}
	Docid *couchbase.DocId
}

func NewViewLocationLessThan(offset interface{}, inclusive bool) *ViewLocation {
	rv := &ViewLocation{
		Key:   offset,
		Docid: &MAX_ID,
	}
	if !inclusive {
		rv.Docid = &MIN_ID
	}
	return rv
}

func NewViewLocationGreatherThan(offset interface{}, inclusive bool) *ViewLocation {
	rv := &ViewLocation{
		Key:   offset,
		Docid: &MIN_ID,
	}
	if !inclusive {
		rv.Docid = &MAX_ID
	}
	return rv
}

func (this *ViewLocation) Compare(that *ViewLocation) int {
	result := ast.CollateJSON(this.Key, that.Key)
	if result != 0 {
		// if keys are the same, compare docis
		if *this.Docid < *that.Docid {
			result = -1
		} else if *this.Docid > *that.Docid {
			result = 1
		} else {
			result = 0
		}
	}

	return result
}

type ViewRange struct {
	Start *ViewLocation
	End   *ViewLocation
}

func (this *ViewRange) Contains(location *ViewLocation) bool {
	rv := this.Start.Compare(location) <= 0 && this.End.Compare(location) >= 0
	return rv
}

func (this *ViewRange) IsSubsetOf(that *ViewRange) bool {
	return that.Contains(this.Start) && that.Contains(this.End)
}

func (this *ViewRange) AsViewQueryOptions() map[string]interface{} {
	return map[string]interface{}{
		"startkey":       this.Start.Key,
		"startkey_docid": *this.Start.Docid,
		"endkey":         this.End.Key,
		"endkey_docid":   *this.End.Docid,
	}
}

func (this *ViewRange) Printable() map[string]interface{} {
	rv := map[string]interface{}{
		"startkey":       this.Start.Key,
		"startkey_docid": *this.Start.Docid,
		"endkey":         this.End.Key,
		"endkey_docid":   *this.End.Docid,
	}

	if this.Start.Docid == &MIN_ID {
		rv["startkey_docid"] = "MIN_DOC_ID"
	}
	if this.End.Docid == &MIN_ID {
		rv["endkey_docid"] = "MIN_DOC_ID"
	}
	if this.Start.Docid == &MAX_ID {
		rv["startkey_docid"] = "MAX_DOC_ID"
	}
	if this.End.Docid == &MAX_ID {
		rv["endkey_docid"] = "MAX_DOC_ID"
	}

	return rv
}

type ViewScanner struct {
	accessPath       *datasource.CouchbaseViewAccessPath
	outputChannel    OutputChannel
	cancelChannel    datasource.CancelChannel
	supportedFactors []ast.BooleanExpression
	ranges           []*ViewRange
}

func NewViewScanner(accessPath *datasource.CouchbaseViewAccessPath) *ViewScanner {
	rv := &ViewScanner{
		accessPath:       accessPath,
		outputChannel:    make(OutputChannel),
		cancelChannel:    make(datasource.CancelChannel),
		supportedFactors: make([]ast.BooleanExpression, 0, 0),
		ranges:           []*ViewRange{&ViewRange{MIN_LOCATION, MAX_LOCATION}},
	}
	return rv
}

func (this *ViewScanner) AddBooleanFactor(factor ast.BooleanExpression) bool {
	// first make sure its sargable
	if !factor.IsSargable() {
		return false
	}

	// now make sure the sarg property matches the index
	sargProperty := factor.GetSargProperty()
	if sargProperty.Path != this.accessPath.Keys()[0] {
		return false
	}

	// calcuate the sarg value
	sargval, err := factor.GetSargValue()
	if err != nil {
		log.Printf("Error determining sarg value: %v", err)
		return false
	}

	// a new array of view ranges to support this boolean factor
	newRanges := make([]*ViewRange, 0, 0)

	switch factor.(type) {
	case *ast.NotEqualToOperator:
		leftRange := &ViewRange{MIN_LOCATION, NewViewLocationLessThan(sargval, false)}
		newRanges = append(newRanges, leftRange)
		rightRange := &ViewRange{NewViewLocationGreatherThan(sargval, false), MAX_LOCATION}
		newRanges = append(newRanges, rightRange)
	case *ast.GreaterThanOperator:
		r := &ViewRange{NewViewLocationGreatherThan(sargval, false), MAX_LOCATION}
		newRanges = append(newRanges, r)
	case *ast.GreaterThanOrEqualOperator:
		r := &ViewRange{NewViewLocationGreatherThan(sargval, true), MAX_LOCATION}
		newRanges = append(newRanges, r)
	case *ast.LessThanOperator:
		r := &ViewRange{MIN_LOCATION, NewViewLocationLessThan(sargval, false)}
		newRanges = append(newRanges, r)
	case *ast.LessThanOrEqualOperator:
		r := &ViewRange{MIN_LOCATION, NewViewLocationLessThan(sargval, true)}
		newRanges = append(newRanges, r)
	case *ast.EqualToOperator:
		r := &ViewRange{NewViewLocationGreatherThan(sargval, true), NewViewLocationLessThan(sargval, true)}
		newRanges = append(newRanges, r)
	}

	// now that we've computed the new ranges, we must merge them with our existing ranges
	for _, r := range newRanges {
		this.MergeRange(r)
	}

	// now that we've updated our rangs, add the factor to the supported factor list
	this.supportedFactors = append(this.supportedFactors, factor)

	return true
}

// the idea here is to combine ranges where we can
// the start state for this is a single range from MIN_LOCATION - MAX_LOCATION
// if you add x < 7
// you should get MIN_LOCATION - 7
// if you then added x > 3
// you should get 3 - 7
// sometimes there are ranges which cannot be merged
// consider x < 3 and x > 7
// in this case we should get 2 ranges
// MIN_LOCATION - 3 and 7 - MAX_LOCATION
// this should also eliminiate some redundant expressions
// consider x < 3 and x < 7
func (this *ViewScanner) MergeRange(newRange *ViewRange) {
	log.Printf("considering new range %v", newRange)
	for _, r := range this.ranges {
		// if the new range is a subset of an existing range
		// this new ragne replaces it
		if newRange.IsSubsetOf(r) {
			log.Printf("new is subset, shrinking")
			r.Start = newRange.Start
			r.End = newRange.End
			return
		}

		if r.IsSubsetOf(newRange) {
			log.Printf("old subset, dropping")
			// drop this factor
			// we already have a stricter range that contains it
			return
		}
	}

	log.Printf("couldn't merge, so add")
	// if we got here, we add it as a new range
	this.ranges = append(this.ranges, newRange)
}

func (this *ViewScanner) GetOutputChannel() OutputChannel {
	return this.outputChannel
}

func (this *ViewScanner) Run() {
	defer close(this.outputChannel)

	for _, r := range this.ranges {
		docChannel := make(datasource.DocumentChannel)
		options := r.AsViewQueryOptions()
		options["reduce"] = false

		go this.accessPath.Scan(docChannel, this.cancelChannel, options)
		for doc := range docChannel {
			this.outputChannel <- doc
		}
	}

}

func (this *ViewScanner) Explain() map[string]interface{} {
	rv := map[string]interface{}{
		"type":           "scan",
		"index":          this.accessPath.Name(),
		"estimated_rows": this.EstimatedRows(),
		"cost":           this.Cost(),
	}

	// ranges can contain some invalid JSON, so rewrite it clean
	printableRanges := make([]map[string]interface{}, 0, len(this.ranges))
	for _, r := range this.ranges {
		printableRanges = append(printableRanges, r.Printable())
	}
	rv["ranges"] = printableRanges

	return rv
}
func (this *ViewScanner) Cancel() {
	close(this.cancelChannel)
}

func (this *ViewScanner) Cost() float64 {
	return (float64(this.EstimatedRows()) / datasource.BATCH_SIZE) * NETWORK_COST
}

func (this *ViewScanner) EstimatedRows() int {
	// if we have no stats at all
	rv := this.accessPath.DataSource().Rows()

	// FIXME only works with 1 level index
	pathStats := this.accessPath.DataSource().PathStats()
	indexKey := this.accessPath.Keys()[0]
	pathStat, ok := pathStats[indexKey]
	if ok {
		// we have path stats for this index
		rv = pathStat.Rows

		// now compute the selectivity factor of all the boolean factors we're supporting
		sf := 1.0
		for _, factor := range this.supportedFactors {
			sf = sf * factor.GetSelectivity(pathStats)
		}

		// multiply the selectivity by the number of rows
		rv = int(float64(rv) * sf)
	} else {
		log.Printf("no path stats for index %v", indexKey)
		log.Printf("%v", pathStats)
	}

	return rv
}

func (this *ViewScanner) TotalCost() float64 { return this.Cost() }

func (this *ViewScanner) Source() Operator {
	return nil
}
