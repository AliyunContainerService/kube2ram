package cr

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

// CreateRepoWebhook invokes the cr.CreateRepoWebhook API synchronously
// api document: https://help.aliyun.com/api/cr/createrepowebhook.html
func (client *Client) CreateRepoWebhook(request *CreateRepoWebhookRequest) (response *CreateRepoWebhookResponse, err error) {
	response = CreateCreateRepoWebhookResponse()
	err = client.DoAction(request, response)
	return
}

// CreateRepoWebhookWithChan invokes the cr.CreateRepoWebhook API asynchronously
// api document: https://help.aliyun.com/api/cr/createrepowebhook.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateRepoWebhookWithChan(request *CreateRepoWebhookRequest) (<-chan *CreateRepoWebhookResponse, <-chan error) {
	responseChan := make(chan *CreateRepoWebhookResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateRepoWebhook(request)
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

// CreateRepoWebhookWithCallback invokes the cr.CreateRepoWebhook API asynchronously
// api document: https://help.aliyun.com/api/cr/createrepowebhook.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateRepoWebhookWithCallback(request *CreateRepoWebhookRequest, callback func(response *CreateRepoWebhookResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateRepoWebhookResponse
		var err error
		defer close(result)
		response, err = client.CreateRepoWebhook(request)
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

// CreateRepoWebhookRequest is the request struct for api CreateRepoWebhook
type CreateRepoWebhookRequest struct {
	*requests.RoaRequest
	RepoNamespace string `position:"Path" name:"RepoNamespace"`
	RepoName      string `position:"Path" name:"RepoName"`
}

// CreateRepoWebhookResponse is the response struct for api CreateRepoWebhook
type CreateRepoWebhookResponse struct {
	*responses.BaseResponse
}

// CreateCreateRepoWebhookRequest creates a request to invoke CreateRepoWebhook API
func CreateCreateRepoWebhookRequest() (request *CreateRepoWebhookRequest) {
	request = &CreateRepoWebhookRequest{
		RoaRequest: &requests.RoaRequest{},
	}
	request.InitWithApiInfo("cr", "2016-06-07", "CreateRepoWebhook", "/repos/[RepoNamespace]/[RepoName]/webhooks", "", "")
	request.Method = requests.PUT
	return
}

// CreateCreateRepoWebhookResponse creates a response to parse from CreateRepoWebhook response
func CreateCreateRepoWebhookResponse() (response *CreateRepoWebhookResponse) {
	response = &CreateRepoWebhookResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
