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
)

type DataSourceManager interface {
	GetDataSource(string) (DataSource, error)
}

type CouchbaseDataSourceManager struct {
	client      couchbase.Client
	pool        couchbase.Pool
	dataSources map[string]*CouchbaseDataSource
}

func NewCouchbaseDataSourceManager(client couchbase.Client) (*CouchbaseDataSourceManager, error) {
	rv := CouchbaseDataSourceManager{
		client:      client,
		dataSources: make(map[string]*CouchbaseDataSource),
	}

	pool, err := client.GetPool("default")
	if err == nil {
		rv.pool = pool
	}

	for k, _ := range pool.BucketMap {
		bucket, err := pool.GetBucket(k)
		if err == nil {
			ds := NewCouchbaseDataSource(bucket)
			rv.dataSources[k] = ds
		}
	}

	log.Printf("Discovered the following datasources: %v", rv.dataSources)

	return &rv, err
}

func (this *CouchbaseDataSourceManager) GetDataSource(name string) (DataSource, error) {
	ds, ok := this.dataSources[name]
	if !ok {
		return nil, fmt.Errorf("No such datasource %v", name)
	}
	return ds, nil
}
