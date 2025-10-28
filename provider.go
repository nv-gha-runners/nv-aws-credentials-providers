/*
 * Copyright (c), NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type credentialFileFormat struct {
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string
	SessionToken    string
	Expiration      *time.Time
}

type StaticWithExpireCredentialsProvider struct {
	credentialsFilePath string
}

func NewStaticWithExpireCredentialsProvider(path string) StaticWithExpireCredentialsProvider {
	return StaticWithExpireCredentialsProvider{credentialsFilePath: path}
}

func (p StaticWithExpireCredentialsProvider) Retrieve(_ context.Context) (aws.Credentials, error) {
	creds := aws.Credentials{}

	fileContent, err := os.ReadFile(p.credentialsFilePath)
	if err != nil {
		return creds, fmt.Errorf("failed to read file at path %v: %w", p.credentialsFilePath, err)
	}

	credsFromFile := credentialFileFormat{}
	err = json.Unmarshal(fileContent, &credsFromFile)
	if err != nil {
		return creds, fmt.Errorf("failed to unmarshal file: %w", err)
	}

	creds.AccessKeyID = credsFromFile.AccessKeyID
	creds.SecretAccessKey = credsFromFile.SecretAccessKey
	creds.SessionToken = credsFromFile.SessionToken

	if credsFromFile.Expiration != nil {
		creds.CanExpire = true
		creds.Expires = *credsFromFile.Expiration
	}

	return creds, nil
}
