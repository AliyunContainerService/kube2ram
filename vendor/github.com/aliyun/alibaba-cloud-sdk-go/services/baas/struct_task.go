package baas

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// Task is a nested struct in baas response
type Task struct {
	TaskId        int    `json:"TaskId" xml:"TaskId"`
	Action        string `json:"Action" xml:"Action"`
	Target        string `json:"Target" xml:"Target"`
	Sender        string `json:"Sender" xml:"Sender"`
	RequestTime   int    `json:"RequestTime" xml:"RequestTime"`
	ResponseTime  string `json:"ResponseTime" xml:"ResponseTime"`
	Handled       bool   `json:"Handled" xml:"Handled"`
	Result        string `json:"Result" xml:"Result"`
	OperationType string `json:"OperationType" xml:"OperationType"`
}
