build:
	GOARCH=amd64 GOOS=linux go build -o artifact/unzipper

deploy: build
	sam package \
		--template-file template.yml \
		--s3-bucket stack-bucket-for-lambda-unzipper \
		--output-template-file sam.yml
	sam deploy \
		--template-file sam.yml \
		--stack-name stack-unzipper-lambda \
		--capabilities CAPABILITY_IAM
