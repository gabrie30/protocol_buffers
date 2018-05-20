## Commands

## Generate protobufs
- $ protoc --go_out=. todo.proto
- $ cat mydb.pd | protoc --decode_raw
- $ hexdump -c mydb.pd

## Generate grpc
- $ protoc -I . todo.proto --go_out=.
- $ protoc -I . todo.proto --go_out=plugins=grpc:.

## Use
cd /cmd/server/ && go build
go run main.go list (fix later but server needs access to db)
go run main.go add some todo to do
