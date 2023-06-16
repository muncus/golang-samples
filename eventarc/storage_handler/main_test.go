// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHelloStorage(t *testing.T) {

	so := storagedata.StorageObjectData{
		Bucket:  "my-example-bucket",
		Name:    "my-object-name",
		Updated: timestamppb.New(time.Now()),
	}
	jsondata, _ := protojson.Marshal(&so)
	ce := cloudevents.NewEvent()
	ce.SetID("sample-id")
	ce.SetSource("//sample/source")
	ce.SetType("google.cloud.storage.object.v1.finalized")
	ce.SetData(*cloudevents.StringOfApplicationJSON(), jsondata)

	w := httptest.NewRecorder()
	r, err := cloudevents.NewHTTPRequestFromEvent(context.Background(), "http://localhost", ce)
	if err != nil {
		t.Fatal(err)
	}
	HelloStorage(w, r)

	want := "my-example-bucket/my-object-name"
	if !strings.Contains(w.Body.String(), want) {
		t.Errorf("%s: got '%s', want '%s'", t.Name(), w.Body.String(), want)
	}
}
