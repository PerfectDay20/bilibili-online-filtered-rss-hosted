.PHONY: build
build:
	GOOS=linux GOARCH=arm64 go build -ldflags='-s' -o bootstrap && zip lambda-handler.zip bootstrap

# upload in the github actions, use env variables as secrets
.PHONY: upload
upload: build
	aws lambda update-function-code --function-name rss --zip-file fileb://lambda-handler.zip

# upload in the local env, use a personal aws profile
.PHONY: upload-local
upload-local: build
	aws --profile personal lambda update-function-code --function-name rss --zip-file fileb://lambda-handler.zip
