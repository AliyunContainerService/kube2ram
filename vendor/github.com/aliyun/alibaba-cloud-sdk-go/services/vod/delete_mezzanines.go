package vod

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

// DeleteMezzanines invokes the vod.DeleteMezzanines API synchronously
// api document: https://help.aliyun.com/api/vod/deletemezzanines.html
func (client *Client) DeleteMezzanines(request *DeleteMezzaninesRequest) (response *DeleteMezzaninesResponse, err error) {
	response = CreateDeleteMezzaninesResponse()
	err = client.DoAction(request, response)
	return
}

// DeleteMezzaninesWithChan invokes the vod.DeleteMezzanines API asynchronously
// api document: https://help.aliyun.com/api/vod/deletemezzanines.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteMezzaninesWithChan(request *DeleteMezzaninesRequest) (<-chan *DeleteMezzaninesResponse, <-chan error) {
	responseChan := make(chan *DeleteMezzaninesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DeleteMezzanines(request)
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

// DeleteMezzaninesWithCallback invokes the vod.DeleteMezzanines API asynchronously
// api document: https://help.aliyun.com/api/vod/deletemezzanines.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DeleteMezzaninesWithCallback(request *DeleteMezzaninesRequest, callback func(response *DeleteMezzaninesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DeleteMezzaninesResponse
		var err error
		defer close(result)
		response, err = client.DeleteMezzanines(request)
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

// DeleteMezzaninesRequest is the request struct for api DeleteMezzanines
type DeleteMezzaninesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	Force                requests.Boolean `position:"Query" name:"Force"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	VideoIds             string           `position:"Query" name:"VideoIds"`
}

// DeleteMezzaninesResponse is the response struct for api DeleteMezzanines
type DeleteMezzaninesResponse struct {
	*responses.BaseResponse
	RequestId            string   `json:"RequestId" xml:"RequestId"`
	NonExistVideoIds     []string `json:"NonExistVideoIds" xml:"NonExistVideoIds"`
	UnRemoveableVideoIds []string `json:"UnRemoveableVideoIds" xml:"UnRemoveableVideoIds"`
}

// CreateDeleteMezzaninesRequest creates a request to invoke DeleteMezzanines API
func CreateDeleteMezzaninesRequest() (request *DeleteMezzaninesRequest) {
	request = &DeleteMezzaninesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("vod", "2017-03-21", "DeleteMezzanines", "vod", "openAPI")
	return
}

// CreateDeleteMezzaninesResponse creates a response to parse from DeleteMezzanines response
func CreateDeleteMezzaninesResponse() (response *DeleteMezzaninesResponse) {
	response = &DeleteMezzaninesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
