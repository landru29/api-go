test:
	go test ./...

swag:
	./swagger-go generate spec -o swagger/swagger.json

install:
	bash ./get-swag.sh
