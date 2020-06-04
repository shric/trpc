package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hekmon/transmissionrpc"
	"github.com/shric/go-trpc/pkg/add"
	"github.com/shric/go-trpc/pkg/list"
)

type dispatch map[string]func(*transmissionrpc.Client, []string)

func main() {
	command := dispatch{
		"list": list.List,
		"add":  add.Add,
	}

	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	addCommand := flag.NewFlagSet("add", flag.ExitOnError)

	if len(os.Args) < 2 {
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	case "add":
		addCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	start := time.Now()

	client := Connect()
	command[os.Args[1]](client, os.Args[2:])
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
