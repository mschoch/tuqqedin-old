//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

type Limit struct {
	source        Operator
	outputChannel OutputChannel
	limit         int
}

func NewLimit(source Operator, limit int) *Limit {
	return &Limit{
		source:        source,
		outputChannel: make(OutputChannel),
		limit:         limit,
	}
}

func (this *Limit) GetOutputChannel() OutputChannel {
	return this.outputChannel
}

func (this *Limit) Run() {
	defer close(this.outputChannel)

	count := 0

	// start the source
	go this.source.Run()
	for row := range this.source.GetOutputChannel() {
		this.outputChannel <- row
		count++
		if count >= this.limit {
			break
		}
	}
}

func (this *Limit) Explain() map[string]interface{} {
	rv := map[string]interface{}{
		"type":           "limit",
		"amount":         this.limit,
		"estimated_rows": this.EstimatedRows(),
		"cost":           this.Cost(),
	}
	if this.Source() != nil {
		rv["source"] = this.Source().Explain()
	}
	return rv
}
func (this *Limit) Cancel() {}

func (this *Limit) Cost() float64 {
	return 0.0
}

func (this *Limit) EstimatedRows() int {
	return this.source.EstimatedRows()
}

func (this *Limit) TotalCost() float64 {
	return this.Cost() + this.source.TotalCost()
}

func (this *Limit) String() string {
	return OperatorToString(this)
}

func (this *Limit) Source() Operator {
	return this.source
}
