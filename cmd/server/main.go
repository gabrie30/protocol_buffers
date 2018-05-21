package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/gabrie30/protocol_buffers/todo"
	"github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
)

type taskServer struct{}
type length int64

const (
	dbPath       = "mydb.pd"
	sizeOfLength = 8
)

var endianness = binary.LittleEndian

func findPath() string {
	absPath, _ := filepath.Abs("../../" + dbPath)
	return absPath
}

func main() {
	srv := grpc.NewServer()
	var tasks taskServer
	todo.RegisterTasksServer(srv, tasks)
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("could not list to :8888: %v", err)
	}

	log.Fatal(srv.Serve(l))
}

func (s taskServer) List(ctx context.Context, void *todo.Void) (*todo.TaskList, error) {

	b, err := ioutil.ReadFile(findPath())
	if err != nil {
		return nil, fmt.Errorf("could not read this file %s: %v", findPath(), err)
	}

	var tasks todo.TaskList
	for {
		if len(b) == 0 {
			return &tasks, nil
		} else if len(b) < sizeOfLength {
			return nil, fmt.Errorf("remaining odd %d bytes len: ", len(b))
		}

		var l length

		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), endianness, &l); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}

		b = b[sizeOfLength:]

		var task todo.Task

		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return nil, fmt.Errorf("could not read task: %v", err)
		}

		b = b[l:]
		tasks.Tasks = append(tasks.Tasks, &task)

		if task.Done {
			fmt.Printf("[X] ")
		} else {
			fmt.Printf("[ ] ")
		}

		fmt.Println(task.Text)
	}
}

func (s taskServer) Add(ctx context.Context, text *todo.Text) (*todo.Task, error) {
	task := &todo.Task{
		Text: text.Text,
		Done: false,
	}

	b, err := proto.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("could not encode task %+v", err)
	}

	f, err := os.OpenFile(findPath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", findPath(), err)
	}

	if err := binary.Write(f, endianness, length(len(b))); err != nil {
		return nil, fmt.Errorf("cound not encode length of message %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file %s: %v", findPath(), err)
	}

	return task, nil
}
