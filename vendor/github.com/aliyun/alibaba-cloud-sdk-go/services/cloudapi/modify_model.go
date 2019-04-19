package cloudapi

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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ModifyModel invokes the cloudapi.ModifyModel API synchronously
// api document: https://help.aliyun.com/api/cloudapi/modifymodel.html
func (client *Client) ModifyModel(request *ModifyModelRequest) (response *ModifyModelResponse, err error) {
	response = CreateModifyModelResponse()
	err = client.DoAction(request, response)
	return
}

// ModifyModelWithChan invokes the cloudapi.ModifyModel API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/modifymodel.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyModelWithChan(request *ModifyModelRequest) (<-chan *ModifyModelResponse, <-chan error) {
	responseChan := make(chan *ModifyModelResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ModifyModel(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ModifyModelWithCallback invokes the cloudapi.ModifyModel API asynchronously
// api document: https://help.aliyun.com/api/cloudapi/modifymodel.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ModifyModelWithCallback(request *ModifyModelRequest, callback func(response *ModifyModelResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ModifyModelResponse
		var err error
		defer close(result)
		response, err = client.ModifyModel(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ModifyModelRequest is the request struct for api ModifyModel
type ModifyModelRequest struct {
	*requests.RpcRequest
	Schema        string `position:"Query" name:"Schema"`
	NewModelName  string `position:"Query" name:"NewModelName"`
	SecurityToken string `position:"Query" name:"SecurityToken"`
	ModelName     string `position:"Query" name:"ModelName"`
	GroupId       string `position:"Query" name:"GroupId"`
	Description   string `position:"Query" name:"Description"`
}

// ModifyModelResponse is the response struct for api ModifyModel
type ModifyModelResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateModifyModelRequest creates a request to invoke ModifyModel API
func CreateModifyModelRequest() (request *ModifyModelRequest) {
	request = &ModifyModelRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudAPI", "2016-07-14", "ModifyModel", "apigateway", "openAPI")
	return
}

// CreateModifyModelResponse creates a response to parse from ModifyModel response
func CreateModifyModelResponse() (response *ModifyModelResponse) {
	response = &ModifyModelResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
