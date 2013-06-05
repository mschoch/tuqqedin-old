//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package datasource

import (
	"fmt"
	"log"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/couchbaselabs/tuqqedin/ast"
	"github.com/couchbaselabs/tuqqedin/stats"
)

type CouchbaseViewAccessPath struct {
	dataSource *CouchbaseDataSource
	ddoc       string
	view       string
	keys       []string
}

func NewCouchbaseViewAccessPath(dataSource *CouchbaseDataSource, ddoc string, view string, keys []string) *CouchbaseViewAccessPath {
	return &CouchbaseViewAccessPath{
		dataSource: dataSource,
		ddoc:       ddoc,
		view:       view,
		keys:       keys,
	}
}

func (this *CouchbaseViewAccessPath) Name() string {
	return fmt.Sprintf("%v%v/_view/%v", DDOC_PREFIX, this.ddoc, this.view)
}

func (this *CouchbaseViewAccessPath) Keys() []string {
	return this.keys
}

func (this *CouchbaseViewAccessPath) DataSource() DataSource {
	return this.dataSource
}

func (this *CouchbaseViewAccessPath) ReturnsAll() bool {
	return false
}

func (this *CouchbaseViewAccessPath) Matches(booleanFactors []ast.BooleanExpression) bool {
	for _, booleanFactor := range booleanFactors {
		switch booleanFactor := booleanFactor.(type) {
		case *ast.OrOperator:
		//FIXME handle partials later
		default:
			rps := booleanFactor.ReferencedProperties()
			if len(rps) == 1 {
				for i, key := range this.keys {
					if i == 0 && key == rps[0].Path {
						// FIXME only matching first index key shortcut
						return true
					}
				}
			}
		}
	}
	return false
}

func (this *CouchbaseViewAccessPath) Scan(output DocumentChannel, cancel CancelChannel, options map[string]interface{}) {
	defer close(output)

	viewRowsChannel := make(chan couchbase.ViewRow)
	go WalkViewInBatches(viewRowsChannel, this.dataSource.bucket, this.ddoc, this.view, options, BATCH_SIZE)
	for row := range viewRowsChannel {
		rowdoc := map[string]interface{}{
			"meta": map[string]interface{}{
				"id": row.ID,
			},
		}
		//		log.Printf("writing doc %v", row)
		output <- rowdoc
	}
}

func (this *CouchbaseViewAccessPath) UpdateStats() {

	log.Printf("Starting UpdateStats for %v", this)

	// only try to address single column stats here
	if len(this.keys) == 1 {

		pathStat := stats.DefaultPathStats(MIN_KEY, MAX_KEY)

		vres, err := this.dataSource.bucket.View(this.ddoc, this.view, map[string]interface{}{"reduce": false, "limit": 0})
		if err != nil {
			log.Printf("Unable to determine cardinality of view, defaulting to MAX")
		} else {
			pathStat.Rows = vres.TotalRows
			pathStat.DistinctValues = vres.TotalRows
		}

		// try to gather deeper stats
		targetCountPerQuantile := pathStat.Rows / pathStat.NumQuantiles()
		options := map[string]interface{}{"group_level": 1}
		viewRowsChannel := make(chan couchbase.ViewRow)
		go WalkViewInBatches(viewRowsChannel, this.dataSource.bucket, this.ddoc, this.view, options, BATCH_SIZE)
		distinctRows := 0
		currentQuantile := stats.QuantileRange{}
		runningCount := 0
		numQuantilesBuilt := 0
		for row := range viewRowsChannel {
			if distinctRows == 0 {
				pathStat.MinValue = row.Key
			}
			if currentQuantile.Count == 0 {
				currentQuantile.Start = row.Key
			}
			pathStat.MaxValue = row.Key
			currentQuantile.End = row.Key
			distinctRows++
			// expect result to be _stats reduce
			switch stats_reduce := row.Value.(type) {
			case map[string]interface{}:
				switch stats_count := stats_reduce["count"].(type) {
				case float64:
					pathStat.MostFrequentValues.Consider(row.Key, stats_count)
					currentQuantile.Count = currentQuantile.Count + int(stats_count)
					runningCount = runningCount + int(stats_count)
				}
			}

			if currentQuantile.Count > targetCountPerQuantile {
				//close out the quantile
				pathStat.Quantiles = append(pathStat.Quantiles, currentQuantile)
				numQuantilesBuilt = numQuantilesBuilt + 1
				// update the target counts (we may have overshot because of a large bin)
				targetCountPerQuantile = (pathStat.Rows - runningCount) / (pathStat.NumQuantiles() - numQuantilesBuilt)
				//empty out a new quantile
				currentQuantile = stats.QuantileRange{}
			}
		}
		// close out the last quantile
		pathStat.Quantiles = append(pathStat.Quantiles, currentQuantile)
		numQuantilesBuilt = numQuantilesBuilt + 1
		pathStat.DistinctValues = distinctRows

		this.dataSource.pathStats[this.keys[0]] = pathStat
		log.Printf("%v", pathStat)
	}

	log.Printf("Finished UpdateStats for %v", this)

}

func (this *CouchbaseViewAccessPath) String() string {
	return fmt.Sprintf("%v", this.Name())
}
