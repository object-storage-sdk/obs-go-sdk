// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of the
// License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.

/**
 * This sample demonstrates how to use custom http.Client to upload multiparts to OBS
 * using the OBS SDK for Go.
 */
package examples

import (
	"fmt"
	"net/http"
	"obs"
	"strings"
)

type SimpleCustumHttpClientSample struct {
	bucketName string
	objectKey  string
	location   string
	obsClient  *obs.ObsClient
}
type OurCustomTransport struct {
	Transport http.RoundTripper
}

func (t *OurCustomTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}
func (t *OurCustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	fmt.Println("do someting")
	return t.transport().RoundTrip(req)
}

func newCustumHttpClientSample(ak, sk, endpoint, bucketName, objectKey, location string) *SimpleCustumHttpClientSample {

	t := &http.Client{
		Transport: &OurCustomTransport{},
	}
	obsClient, err := obs.New(ak, sk, endpoint, obs.WithHttpClient(t))
	if err != nil {
		panic(err)
	}
	return &SimpleCustumHttpClientSample{obsClient: obsClient, bucketName: bucketName, objectKey: objectKey, location: location}
}

func (sample SimpleCustumHttpClientSample) CreateBucket() {
	input := &obs.CreateBucketInput{}
	input.Bucket = sample.bucketName
	input.Location = sample.location
	_, err := sample.obsClient.CreateBucket(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create bucket:%s successfully!\n", sample.bucketName)
	fmt.Println()
}

func (sample SimpleCustumHttpClientSample) InitiateMultipartUpload() string {
	input := &obs.InitiateMultipartUploadInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	output, err := sample.obsClient.InitiateMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	return output.UploadId
}

func (sample SimpleCustumHttpClientSample) UploadPart(uploadId string) (string, int) {
	input := &obs.UploadPartInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.UploadId = uploadId
	input.PartNumber = 1
	input.Body = strings.NewReader("Hello OBS")
	output, err := sample.obsClient.UploadPart(input)
	if err != nil {
		panic(err)
	}
	return output.ETag, output.PartNumber
}

func (sample SimpleCustumHttpClientSample) CompleteMultipartUpload(uploadId, etag string, partNumber int) {
	input := &obs.CompleteMultipartUploadInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.UploadId = uploadId
	input.Parts = []obs.Part{
		obs.Part{PartNumber: partNumber, ETag: etag},
	}
	_, err := sample.obsClient.CompleteMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload object %s successfully!\n", sample.objectKey)
}

func RunCustumHttpClientSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
		location   = "yourbucketlocation"
	)

	sample := newCustumHttpClientSample(ak, sk, endpoint, bucketName, objectKey, location)

	fmt.Println("Create a new bucket for demo")
	sample.CreateBucket()

	// Step 1: initiate multipart upload
	fmt.Println("Step 1: initiate multipart upload")
	uploadId := sample.InitiateMultipartUpload()

	// Step 2: upload a part
	fmt.Println("Step 2: upload a part")

	etag, partNumber := sample.UploadPart(uploadId)

	// Step 3: complete multipart upload
	fmt.Println("Step 3: complete multipart upload")
	sample.CompleteMultipartUpload(uploadId, etag, partNumber)

}
