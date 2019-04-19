package cdn

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

// SetForceRedirectConfig invokes the cdn.SetForceRedirectConfig API synchronously
// api document: https://help.aliyun.com/api/cdn/setforceredirectconfig.html
func (client *Client) SetForceRedirectConfig(request *SetForceRedirectConfigRequest) (response *SetForceRedirectConfigResponse, err error) {
	response = CreateSetForceRedirectConfigResponse()
	err = client.DoAction(request, response)
	return
}

// SetForceRedirectConfigWithChan invokes the cdn.SetForceRedirectConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/setforceredirectconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetForceRedirectConfigWithChan(request *SetForceRedirectConfigRequest) (<-chan *SetForceRedirectConfigResponse, <-chan error) {
	responseChan := make(chan *SetForceRedirectConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.SetForceRedirectConfig(request)
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

// SetForceRedirectConfigWithCallback invokes the cdn.SetForceRedirectConfig API asynchronously
// api document: https://help.aliyun.com/api/cdn/setforceredirectconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) SetForceRedirectConfigWithCallback(request *SetForceRedirectConfigRequest, callback func(response *SetForceRedirectConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *SetForceRedirectConfigResponse
		var err error
		defer close(result)
		response, err = client.SetForceRedirectConfig(request)
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

// SetForceRedirectConfigRequest is the request struct for api SetForceRedirectConfig
type SetForceRedirectConfigRequest struct {
	*requests.RpcRequest
	SecurityToken string           `position:"Query" name:"SecurityToken"`
	DomainName    string           `position:"Query" name:"DomainName"`
	RedirectType  string           `position:"Query" name:"RedirectType"`
	OwnerId       requests.Integer `position:"Query" name:"OwnerId"`
}

// SetForceRedirectConfigResponse is the response struct for api SetForceRedirectConfig
type SetForceRedirectConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateSetForceRedirectConfigRequest creates a request to invoke SetForceRedirectConfig API
func CreateSetForceRedirectConfigRequest() (request *SetForceRedirectConfigRequest) {
	request = &SetForceRedirectConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Cdn", "2018-05-10", "SetForceRedirectConfig", "", "")
	return
}

// CreateSetForceRedirectConfigResponse creates a response to parse from SetForceRedirectConfig response
func CreateSetForceRedirectConfigResponse() (response *SetForceRedirectConfigResponse) {
	response = &SetForceRedirectConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
