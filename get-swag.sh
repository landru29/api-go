#!/bin/bash

latestv=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | jq -r .tag_name)
curl -o ./swagger-go -L'#' https://github.com/go-swagger/go-swagger/releases/download/$latestv/swagger_$(echo `uname`|tr '[:upper:]' '[:lower:]')_amd64
chmod +x ./swagger-go
