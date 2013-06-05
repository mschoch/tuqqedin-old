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
	//	"log"
	"math"
	"strings"

	"github.com/couchbaselabs/go-couchbase"
	"github.com/couchbaselabs/tuqqedin/stats"
)

const VIEW_NAME_PREFIX = "index_by_"
const DEV_DDOC_PREFIX = "dev_"
const DDOC_PREFIX = "_design/"
const BATCH_SIZE = 1000

type CouchbaseDataSource struct {
	bucket    *couchbase.Bucket
	rows      int
	pathStats map[string]stats.PathStatistics
	alldocs   AccessPath
	views     map[string]AccessPath
}

func NewCouchbaseDataSource(bucket *couchbase.Bucket) *CouchbaseDataSource {
	rv := &CouchbaseDataSource{
		bucket:    bucket,
		rows:      math.MaxInt32,
		pathStats: make(map[string]stats.PathStatistics),
		alldocs:   nil,
		views:     make(map[string]AccessPath),
	}

	// do this syncrhonously at startup (should not be too expensive)
	rv.UpdateAccessPaths()

	// then asynchronously retrieve some stats
	go rv.UpdateStats()

	return rv
}

func (this *CouchbaseDataSource) Name() string {
	return this.bucket.Name
}

func (this *CouchbaseDataSource) Rows() int {
	return this.rows
}

func (this *CouchbaseDataSource) PathStats() map[string]stats.PathStatistics {
	return this.pathStats
}

func (this *CouchbaseDataSource) UpdateStats() {
	this.alldocs.UpdateStats()

	for _, view := range this.views {
		view.UpdateStats()
	}
}

func (this *CouchbaseDataSource) AccessPaths() []AccessPath {
	rv := make([]AccessPath, 0, len(this.views)+1)
	rv = append(rv, this.alldocs)
	for _, view := range this.views {
		rv = append(rv, view)
	}
	return rv
}

func (this *CouchbaseDataSource) UpdateAccessPaths() {

	if this.alldocs == nil {
		this.alldocs = NewCouchbaseAllDocsAccessPath(this)
	}

	ddocs := this.getProductionDesignDocuments()
	for _, ddoc := range ddocs {
		for name, _ := range ddoc.Json.Views {

			if strings.HasPrefix(name, VIEW_NAME_PREFIX) {
				index_keys_string := name[len(VIEW_NAME_PREFIX):]
				index_keys := strings.Split(index_keys_string, "_")

				viewAccessPath := NewCouchbaseViewAccessPath(this, designDocName(ddoc), name, index_keys)
				this.views[viewAccessPath.Name()] = viewAccessPath
			}
		}
	}

}

func (this *CouchbaseDataSource) String() string {
	return fmt.Sprintf("AccessPaths: %v", this.AccessPaths())
}

func (this *CouchbaseDataSource) getProductionDesignDocuments() []couchbase.DDoc {
	rv := make([]couchbase.DDoc, 0, 0)

	ddocs, err := this.bucket.GetDDocs()
	if err == nil {
		for _, ddocrow := range ddocs.Rows {
			if !strings.HasPrefix(designDocName(ddocrow.DDoc), DEV_DDOC_PREFIX) {
				rv = append(rv, ddocrow.DDoc)
			}
		}
	}
	return rv
}

func (this *CouchbaseDataSource) Fetch(docID string) (interface{}, error) {
	var rv map[string]interface{}
	err := this.bucket.Get(docID, &rv)
	return rv, err
}

func designDocName(ddoc couchbase.DDoc) string {
	rv := ""
	switch ddocName := ddoc.Meta["id"].(type) {
	case string:
		rv = ddocName
		if strings.HasPrefix(ddocName, DDOC_PREFIX) {
			rv = ddocName[len(DDOC_PREFIX):]
		}
	}
	return rv
}
