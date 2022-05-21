#!/bin/bash

go get github.com/99designs/gqlgen@v0.17.5
go run github.com/99designs/gqlgen generate
sed -i -r 's/`(json:(".*"))/`\1 bson:\2/g' graph/model/models_gen.go
