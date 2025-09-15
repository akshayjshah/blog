---
title: "Conditional S3 writes in Go"
description: "How to issue conditional writes to S3 with aws-sdk-go-v2 and test them with MinIO."
created: 2025-09-13
---

S3 [announced support for conditional writes][announcement] in August of 2024.
Conditional writes allow distributed systems to safely read objects, modify
them, and write them back to S3 without any additional dependencies. This
*greatly* simplifies many read-heavy systems.

However, AWS's documentation for this feature --- especially in Go --- is
terrible. In this post, I'll show you how to issue conditional writes with v2 of
AWS's Go SDK. I'll also show you how to write integration tests with
`testcontainers` and MinIO.

If you'd rather jump straight to Github, all the code in this post is available
in [`akshayjshah/conditionalwrite`][github].

## Issuing conditional writes

Issuing a conditional write is as simple as setting the `If-None-Match` or
`If-Match` HTTP headers. With a small single-object client type:

```go
package conditionalwrite

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type ETag string

const None ETag = ""

type Client struct {
	client *s3.Client
	bucket string
	key    string
}

func (c *Client) Set(
 	ctx context.Context,
 	r io.Reader,
 	previous ETag) (ETag, error) {

	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.key),
		Body:   r,
	}
	if previous == "" {
		input.IfNoneMatch = aws.String("*")
	} else {
		input.IfMatch = aws.String(string(previous))
	}
	res, err := c.client.PutObject(ctx, input)
	if err != nil {
		return "", err
	}
	if res == nil || res.ETag == nil {
		return "", errors.New("no ETag")
	}
	return ETag(*res.ETag), nil
}
```

When clients are in a read-modify-write loop (usually called "optimistic
concurrency control"), it's important to distinguish concurrency control errors.
To do this in Go, we must reach down into the [Smithy][] package:

```go
func IsPreconditionFailed(err error) bool {
	return getSmithyCode(err) == "PreconditionFailed"
}

func getSmithyCode(err error) string {
	if err == nil {
		return ""
	}
	var e smithy.APIError
	if errors.As(err, &e) {
		return e.ErrorCode()
	}
	return ""
}
```

Of course, we'll also need a way to construct a `Client`. I found AWS's
authentication packages difficult to grok --- the documentation rightly focuses
on production use cases, but I wanted to work with local object storage.

```go
func NewClient(endpoint, user, pw, region, bucket, key string) *Client {
	c := s3.New(s3.Options{
		Region:       region,
		BaseEndpoint: aws.String(endpoint),
		DefaultsMode: aws.DefaultsModeStandard,
		Credentials: credentials.NewStaticCredentialsProvider(
			user,
			pw,
			"", /* session */
		),
		UsePathStyle:               true,
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenSupported,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenSupported,
		HTTPClient: &http.Client{
			Transport: &http.Transport{},
		},
	})
	return &Client{client: c, bucket: bucket, key: key}
}
```

For tests, it's also nice to have a method to create buckets:

```go
func (c *Client) CreateBucket(ctx context.Context) error {
	_, err := c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(c.bucket),
	})
	if getSmithyCode(err) == "BucketAlreadyOwnedByYou" {
		return nil
	}
	return err
}
```

## Integration tests with MinIO

I don't like mocking complex dependencies like S3, so I'd prefer to run an
S3-compatible object store in my tests. [MinIO][] is widely used and has an
excellent [`testcontainers`][testcontainers] module, so it's easy to integrate
and allows each test to have an isolated object store. This does mean that tests
require Docker, but I'm happy with that tradeoff.

Here's a simple test that creates an object, updates it successfully, and then
tries to update it with an outdated ETag:

```go
package conditionalwrite

import (
	"fmt"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/minio"
)

func TestConditionalWrite(t *testing.T) {
	// Requires a running (or socket-activated) Docker daemon.
	const user, password = "admin", "password"
	mc, err := minio.Run(
		t.Context(),
		"minio/minio:RELEASE.2025-07-23T15-54-02Z",
		minio.WithUsername(user),
		minio.WithPassword(password),
	)
	if err != nil {
		t.Fatalf("start MinIO container: %v", err)
	}
	addr, err := mc.ConnectionString(t.Context())
	if err != nil {
		t.Fatalf("get MinIO connection string: %v", err)
	}

	c := NewClient(
		fmt.Sprintf("http://%s", addr), // endpoint
		user,
		password,
		"us-east-1", // region
		"test",      // bucket
		"text.txt",  // key
	)

	err = c.CreateBucket(t.Context())
	if err != nil {
		t.Fatalf("create bucket failed: %v", err)
	}

	etag, err := c.Set(t.Context(), strings.NewReader("one"), None)
	if err != nil {
		t.Fatalf("initial write failed: %v", err)
	}

	_, err = c.Set(t.Context(), strings.NewReader("two"), etag)
	if err != nil {
		t.Fatalf("overwrite with correct ETag failed: %v", err)
	}

	_, err = c.Set(t.Context(), strings.NewReader("three"), etag)
	if err == nil {
		t.Fatal("overwrite with incorrect ETag succeeded")
	}
	if !IsPreconditionFailed(err) {
		t.Fatalf("expected PreconditionFailed error, got %v", err)
	}
}
```

The overhead of starting a MinIO container makes this test slow enough that I'd
consider skipping it when `testing.Short()` is set.

## An aside on generated clients

After 19 years, S3's API has grown quite a bit: its [Smithy model][smithy-s3] is
a **forty thousand** line JSON file. Because there's no distinction between
commonly-used and long-tail endpoints, the generated Go package is frustratingly
enormous and cumbersome. I'd never let types from this package leak into
production code. Instead, I'd write a wrapper and enforce its usage with a
custom linter.

That said, I'm glad that S3 finally supports conditional writes. Optimistic
concurrency control with `If-Match` is dead simple, doesn't add any moving
parts, and is efficient enough for many read-dominated workloads.

All the code in this post is [available on Github][github].

[announcement]: https://aws.amazon.com/about-aws/whats-new/2024/08/amazon-s3-conditional-writes/
[github]: https://github.com/akshayjshah/conditionalwrite
[MinIO]: https://www.min.io/
[smithy-s3]: https://raw.githubusercontent.com/aws/api-models-aws/9c9dd620e2541b82f34ac5d52d73625b753f80a8/models/s3/service/2006-03-01/s3-2006-03-01.json
[Smithy]: https://smithy.io/index.html
[testcontainers]: https://testcontainers.com/ 
