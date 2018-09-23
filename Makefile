build:
	GOARCH=amd64 GOOS=linux go build -o artifact/unzipper
.PHONY: build

deploy: build
	sam package \
		--template-file template.yml \
		--s3-bucket stack-bucket-for-lambda-unzipper \
		--output-template-file sam.yml
	sam deploy \
		--template-file sam.yml \
		--stack-name stack-unzipper-lambda \
		--capabilities CAPABILITY_IAM
.PHONY: deploy

delete:
	aws s3 rm s3://zipped-artifact --recursive
	aws s3 rm s3://unzipped-artifact --recursive
	aws cloudformation delete-stack --stack-name stack-unzipper-lambda
	aws s3 rm s3://stack-bucket-for-lambda-unzipper --recursive
	aws s3 rb s3://stack-bucket-for-lambda-unzipper
.PHONY: delete

test:
	go test ./...
.PHONY: test
