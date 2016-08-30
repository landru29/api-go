# API


## Develop

```
go get github.com/yvasiyarov/swagger
cd $GOPATH/src/github.com/yvasiyarov/swagger
go build
go install
cd $GOPATH/src/github.com/landru29/api-go
$GOPATH/bin/swagger -apiPackage="github.com/landru29/api-go" -mainApiFile="github.com/landru29/api-go/api.go" -format="swagger"
```
