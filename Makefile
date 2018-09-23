build:
	GOARCH=amd64 GOOS=linux go build -o artifact/unzipper

deploy: build
	sam package \
		--template-file template.yml \
		--s3-bucket lambda-unzipper \
		--output-template-file sam.yml
	sam deploy \
		--template-file sam.yml \
		--stack-name stack-unzipper-lambda \
		--capabilities CAPABILITY_IAM

delete:
	aws s3 rm s3://zipped-artifact --recursive
	aws s3 rm s3://unzipped-artifact --recursive
	aws cloudformation delete-stack --stack-name stack-unzipper-lambda
	aws s3 rm s3://lambda-unzipper --recursive
    aws s3 rb s3://lambda-unzipper
