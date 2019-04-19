package imm

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

// AudiosItem is a nested struct in imm response
type AudiosItem struct {
	AudioDuration        float64          `json:"AudioDuration" xml:"AudioDuration"`
	AudioUri             string           `json:"AudioUri" xml:"AudioUri"`
	AudioRate            int              `json:"AudioRate" xml:"AudioRate"`
	SourceType           string           `json:"SourceType" xml:"SourceType"`
	ModifyTime           string           `json:"ModifyTime" xml:"ModifyTime"`
	FileSize             int              `json:"FileSize" xml:"FileSize"`
	SourceUri            string           `json:"SourceUri" xml:"SourceUri"`
	CreateTime           string           `json:"CreateTime" xml:"CreateTime"`
	RemarksA             string           `json:"RemarksA" xml:"RemarksA"`
	RemarksB             string           `json:"RemarksB" xml:"RemarksB"`
	AudioTextsStatus     string           `json:"AudioTextsStatus" xml:"AudioTextsStatus"`
	AudioTextsModifyTime string           `json:"AudioTextsModifyTime" xml:"AudioTextsModifyTime"`
	ProcessModifyTime    string           `json:"ProcessModifyTime" xml:"ProcessModifyTime"`
	ProcessStatus        string           `json:"ProcessStatus" xml:"ProcessStatus"`
	SourcePosition       string           `json:"SourcePosition" xml:"SourcePosition"`
	AudioFormat          string           `json:"AudioFormat" xml:"AudioFormat"`
	AudioTextsFailReason string           `json:"AudioTextsFailReason" xml:"AudioTextsFailReason"`
	ProcessFailReason    string           `json:"ProcessFailReason" xml:"ProcessFailReason"`
	RemarksC             string           `json:"RemarksC" xml:"RemarksC"`
	RemarksD             string           `json:"RemarksD" xml:"RemarksD"`
	ExternalId           string           `json:"ExternalId" xml:"ExternalId"`
	AudioTexts           []AudioTextsItem `json:"AudioTexts" xml:"AudioTexts"`
}
