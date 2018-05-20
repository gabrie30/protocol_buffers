package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gabrie30/protocol_buffers/todo"
	grpc "google.golang.org/grpc"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: list or add")
		os.Exit(1)
	}
	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to backend: %v", err)
	}

	client := todo.NewTasksClient(conn)

	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list(context.Background(), client)
	case "add":
		err = add(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const (
	dbPath       = "mydb.pd"
	sizeOfLength = 8
)

var endianness = binary.LittleEndian

func add(ctx context.Context, client todo.TasksClient, text string) error {
	_, err := client.Add(ctx, &todo.Text{Text: text})

	if err != nil {
		return fmt.Errorf("Could not create task %v", err)
	}

	fmt.Println("task added successfully")
	return nil
}

func list(ctx context.Context, client todo.TasksClient) error {
	l, err := client.List(ctx, &todo.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch tasks: %v", err)
	}

	for _, t := range l.Tasks {
		if t.Done {
			fmt.Printf("[X] ")
		} else {
			fmt.Printf("[ ] ")
		}

		fmt.Println(t.Text)
	}

	return nil
}
