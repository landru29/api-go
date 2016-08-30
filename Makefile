test:
	go test ./...

swag:
	./swagger-go generate spec -o swagger/dist/swagger.json

install:
	bash ./get-swag.sh
