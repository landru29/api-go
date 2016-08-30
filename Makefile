test:
	go test ./...

swag:
	swagger generate spec -o swagger/dist/swagger.json
