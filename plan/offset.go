//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

type Offset struct {
	source        Operator
	outputChannel OutputChannel
	offset        int
}

func NewOffset(source Operator, offset int) *Offset {
	return &Offset{
		source:        source,
		outputChannel: make(OutputChannel),
		offset:        offset,
	}
}

func (this *Offset) GetOutputChannel() OutputChannel {
	return this.outputChannel
}

func (this *Offset) Run() {
	defer close(this.outputChannel)

	count := 0

	// start the source
	go this.source.Run()
	for row := range this.source.GetOutputChannel() {
		count++
		if count <= this.offset {
			continue
		}
		this.outputChannel <- row
	}
}

func (this *Offset) Explain() map[string]interface{} {
	rv := map[string]interface{}{
		"type":           "offset",
		"amount":         this.offset,
		"estimated_rows": this.EstimatedRows(),
		"cost":           this.Cost(),
	}
	if this.Source() != nil {
		rv["source"] = this.Source().Explain()
	}
	return rv
}
func (this *Offset) Cancel() {}

func (this *Offset) Cost() float64 {
	return 0.0
}

func (this *Offset) EstimatedRows() int {
	return this.source.EstimatedRows()
}

func (this *Offset) TotalCost() float64 {
	return this.Cost() + this.source.TotalCost()
}

func (this *Offset) String() string {
	return OperatorToString(this)
}

func (this *Offset) Source() Operator {
	return this.source
}
