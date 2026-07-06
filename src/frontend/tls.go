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
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc/credentials"
)

func loadClientCredentials(address string) (credentials.TransportCredentials, error) {
	caPath := os.Getenv("GRPC_TLS_CA")
	if caPath == "" {
		return nil, fmt.Errorf("environment variable %q not set", "GRPC_TLS_CA")
	}
	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return nil, fmt.Errorf("read gRPC CA certificate: %w", err)
	}

	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("parse gRPC CA certificate %q", caPath)
	}

	serverName, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("parse gRPC service address %q: %w", address, err)
	}

	return credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS13,
		RootCAs:    roots,
		ServerName: serverName,
	}), nil
}
