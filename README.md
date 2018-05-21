# Commands

## Generate protobufs
```bash
$ protoc --go_out=. todo.proto
$ cat mydb.pd | protoc --decode_raw
$ hexdump -c mydb.pd
```

## Generate grpc

```bash
$ protoc -I . todo.proto --go_out=.
$ protoc -I . todo.proto --go_out=plugins=grpc:.
```

## Use

```bash
cd ./cmd/server/ && go run main.go
cd ./cmd/todo/ && go run main.go add first todo
go run main.go list
```
