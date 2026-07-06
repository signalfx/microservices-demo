// Copyright 2026 Splunk Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/pem"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoadClientCredentialsVerifiesTrustedServer(t *testing.T) {
	testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	testServer.EnableHTTP2 = true
	testServer.StartTLS()
	t.Cleanup(testServer.Close)

	caPath := t.TempDir() + "/ca.crt"
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: testServer.Certificate().Raw})
	if err := os.WriteFile(caPath, caPEM, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GRPC_TLS_CA", caPath)

	transportCredentials, err := loadClientCredentials("example.com:443")
	if err != nil {
		t.Fatal(err)
	}
	rawConnection, err := net.Dial("tcp", strings.TrimPrefix(testServer.URL, "https://"))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	secureConnection, _, err := transportCredentials.ClientHandshake(ctx, "example.com:443", rawConnection)
	if err != nil {
		t.Fatal(err)
	}
	secureConnection.Close()
}
