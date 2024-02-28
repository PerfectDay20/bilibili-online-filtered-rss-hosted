.PHONY: build
build:
	GOOS=linux GOARCH=arm64 go build -o bootstrap && zip lambda-handler.zip bootstrap

.PHONY: upload
upload: build
	aws --profile personal lambda update-function-code --function-name rss --zip-file fileb://lambda-handler.zip
