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
	"crypto/tls"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

func loadServerCredentials() (credentials.TransportCredentials, error) {
	certPath := os.Getenv("GRPC_TLS_CERT")
	if certPath == "" {
		return nil, fmt.Errorf("environment variable %q not set", "GRPC_TLS_CERT")
	}
	keyPath := os.Getenv("GRPC_TLS_KEY")
	if keyPath == "" {
		return nil, fmt.Errorf("environment variable %q not set", "GRPC_TLS_KEY")
	}

	certificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("load gRPC server certificate: %w", err)
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS13,
	}), nil
}
