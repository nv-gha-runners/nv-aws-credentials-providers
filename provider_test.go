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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

var _ aws.CredentialsProvider = (*StaticWithExpireCredentialsProvider)(nil)

func TestRetrieveInvalidFilePath(t *testing.T) {
	p := NewStaticWithExpireCredentialsProvider("/invalid/path")
	_, err := p.Retrieve(context.TODO())
	assert.Error(t, err)
}

func TestRetrieve(t *testing.T) {
	expectedExpirationString := "2025-10-28T18:05:26Z"
	expectedExpiration, err := time.Parse(time.RFC3339, expectedExpirationString)
	assert.NoError(t, err)

	expectedAccesKeyId := "accessKey"
	expectedSecretKey := "secretKey"
	expectedSessionToken := "sessionToken"

	testCases := []struct {
		fileContent             string
		shouldFail              bool
		shouldExpire            bool
		expectedAccesKeyId      string
		expectedSecretAccessKey string
		expectedSessionToken    string
		expectedExpiration      time.Time
	}{
		{
			fileContent: fmt.Sprintf(`{
				"Version": 1,
				"AccessKeyId": "%s",
				"SecretAccessKey": "%s",
				"SessionToken": "%s",
				"Expiration": "%s"}`,
				expectedAccesKeyId, expectedSecretKey, expectedSessionToken, expectedExpirationString),
			shouldFail:              false,
			shouldExpire:            true,
			expectedAccesKeyId:      expectedAccesKeyId,
			expectedSecretAccessKey: expectedSecretKey,
			expectedSessionToken:    expectedSessionToken,
			expectedExpiration:      expectedExpiration,
		},
		{
			fileContent: fmt.Sprintf(`{
				"Version": 1,
				"AccessKeyId": "%s",
				"SecretAccessKey": "%s"
				}`, expectedAccesKeyId, expectedSecretKey),
			shouldFail:              false,
			shouldExpire:            false,
			expectedAccesKeyId:      expectedAccesKeyId,
			expectedSecretAccessKey: expectedSecretKey,
			expectedSessionToken:    "",
			expectedExpiration:      time.Time{},
		},
		{
			fileContent:             "invalid",
			shouldFail:              true,
			shouldExpire:            false,
			expectedAccesKeyId:      "",
			expectedSecretAccessKey: "",
			expectedSessionToken:    "",
			expectedExpiration:      time.Time{},
		},
	}

	for _, tc := range testCases {
		file, err := os.CreateTemp("", "test")
		assert.NoError(t, err)
		assert.NotNil(t, file)

		_, err = file.WriteString(tc.fileContent)
		assert.NoError(t, err)

		p := NewStaticWithExpireCredentialsProvider(file.Name())
		creds, err := p.Retrieve(context.TODO())

		if tc.shouldFail {
			assert.Error(t, err)
			continue
		}

		assert.NoError(t, err)
		assert.Equal(t, tc.expectedAccesKeyId, creds.AccessKeyID)
		assert.Equal(t, tc.expectedSecretAccessKey, creds.SecretAccessKey)
		assert.Equal(t, tc.expectedSessionToken, creds.SessionToken)
		assert.Equal(t, tc.expectedExpiration, creds.Expires)
		assert.Equal(t, tc.shouldExpire, creds.CanExpire)

		err = os.Remove(file.Name())
		assert.NoError(t, err)
	}
}
