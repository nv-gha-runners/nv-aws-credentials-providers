# nv-aws-credentials-providers

This repository contains a set of AWS Credentials Providers for the AWS SDK for Go v2.

## StaticWithExpireCredentialsProvider

This provider is getting AWS Credentials from a file, as the default AWS file
provider, but also handles expiration of credentials. This is done by using the
file format expected by the [Credential Process Provider](https://docs.aws.amazon.com/sdkref/latest/guide/feature-process-credentials.html#feature-process-credentials-output).

A valid credentials file would like this:
```json
{
    "Version": 1,
    "AccessKeyId": "an AWS access key",
    "SecretAccessKey": "your AWS secret access key",
    "SessionToken": "the AWS session token for temporary credentials",
    "Expiration": "RFC3339 timestamp for when the credentials expire"
}
```
Only the `AccessKeyId` and `SecretAccessKey` fields are required.
This file should be updated by an external tool when the credentials are expiring.

Here is a simple example on how to use this provider:
```go
package main

import (
        "context"
        "log"

        "github.com/aws/aws-sdk-go-v2/config"
        provider "github.com/nv-gha-runners/nv-aws-credentials-providers"
)

func main() {
        p := provider.NewStaticWithExpireCredentialsProvider("/path-to-creds-file")
        awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(p))
        if err != nil {
                log.Fatalf("failed to create AWS config: %v\n", err)
        }

        [...]
}
```
