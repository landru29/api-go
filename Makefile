test:
	go test ./...

swag:
	swaggo --output swagger/swagger.json

install:
	bash ./get-swag.sh

build: swag
	go build
