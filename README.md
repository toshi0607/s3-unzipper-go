s3-unzipper-go
====

## Description
unzip uploaded zip file to another S3 bucket via AWS Lambda in Go

## setup

install libraries

```
$ dep ensure
```

set environment variables for test

```
$ export SRC_BUCKET=zipped-artifact-dev
$ export DEST_BUCKET=unzipped-artifact-dev
```

## Local
You can test a behavior on test (`main_test.go`).

In the test, `setup` prepares 2 real S3 buckets because SAM local doesn't support local emulation of an S3.

One is for an even source that triggers an AWS Lambda and another is for a destination of unzipped artifacts.

Because S3 buckets created at the test are deleted on every test execution, idempotency is guaranteed.

## Production

### prerequisites

You have to prepare credentials with proper policies.

* AWSLambdaFullAccess (should be limited)
* AmazonS3FullAccess (should be limited)
* CloudWatchLogsFullAccess (should be limited)
* AWSXrayFullAccess (should be limited)

In addition to the above, I added group policy for AWS CloudFormation. It's because `sam deploy` command is alias of `aws cloudformation deploy`.

Also, `sam package` command generates CloudFormation template from *template.yml*.

### deploy

```
$ make deploy
```
