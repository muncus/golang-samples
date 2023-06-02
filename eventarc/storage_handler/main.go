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

// [START eventarc_storage_cloudevent_handler]

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/googleapis/google-cloudevents-go/cloud/storagedata"
	"google.golang.org/protobuf/encoding/protojson"
)

// HelloStorage receives and processes a CloudEvent containing a StorageObjectData
func HelloStorage(w http.ResponseWriter, r *http.Request) {
	ce, err := cloudevents.NewEventFromHTTPRequest(r)
	if err != nil {
		log.Printf("failed to parse CloudEvent: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	var storageevent storagedata.StorageObjectData
	err = protojson.Unmarshal(ce.Data(), &storageevent)
	if err != nil {
		log.Printf("failed to unmarshal: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	// TODO: consider using .SelfLink here?
	s := fmt.Sprintf("Cloud Storage object changed: %s updated at %s",
		path.Join(storageevent.GetBucket(), storageevent.GetName()),
		storageevent.Updated.AsTime().UTC())
	log.Printf(s)
	fmt.Fprintln(w, s)
}

// [END eventarc_storage_cloudevent_handler]
// [START eventarc_storage_cloudevent_server]

func main() {
	http.HandleFunc("/", HelloStorage)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start HTTP server.
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// [END eventarc_storage_cloudevent_server]
