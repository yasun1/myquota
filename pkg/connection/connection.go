package connection

/*
Copyright (c) 2018 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	client "github.com/openshift-online/ocm-sdk-go"
)

const (
	tokenURL       = "https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token"
	clientID       = "cloud-services"
	clientSecret   = ""
	skipAuth       = true
	integration    = false
	healthcheckURL = "http://localhost:8083"
)

func gatewayURL() (url string) {
	ocmEnv := os.Getenv("OCM_ENV")
	switch ocmEnv {
	case "production":
		url = "https://api.openshift.com"
	case "staging":
		url = "https://api.stage.openshift.com"
	case "integration":
		url = "https://api.integration.openshift.com"
	default:
		url = "https://api.stage.openshift.com"
	}

	return url
}

// SuperAdmin
var (
	superAdminUserToken  = os.Getenv("SUPER_ADMIN_USER_TOKEN")
	SuperAdminConnection = createConnectionWithToken(superAdminUserToken)
)

var (
	// Create a logger:
	logger = createLogger()
)

func createConnectionWithToken(token string) *client.Connection {
	gatewayURL := gatewayURL()

	if token == "" {
		fmt.Println("[WARNING]: Token shouldn't be empty")
	}

	// Create the connection:
	connection, err := client.NewConnectionBuilder().
		Logger(logger).
		Insecure(true).
		TokenURL(tokenURL).
		URL(gatewayURL).
		Client(clientID, clientSecret).
		Tokens(token).
		Build()
	if err != nil {
		fmt.Printf("ERROR occurred when create connection with token: %s!! %s\n", token, err)
	}
	return connection

}

func createLogger() client.Logger {
	debugMode := false
	if os.Getenv("OCM_Debug_Mode") == "true" {
		debugMode = true
	}
	logger, _ := client.NewStdLoggerBuilder().
		Streams(GinkgoWriter, GinkgoWriter).
		Debug(debugMode).
		Build()
	// ExpectWithOffset(1, err).ToNot(HaveOccurred())
	// logger, _ := client.NewGlogLoggerBuilder().DebugV(3).InfoV(0).Build()

	return logger
}

func checkSkipAuth() {
	if skipAuth {
		Skip("Test skipped due to authorization issue")
	}
}
