package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hekmon/transmissionrpc"
	"github.com/jessevdk/go-flags"
	"github.com/shric/go-trpc/pkg/add"
	"github.com/shric/go-trpc/pkg/list"
	"github.com/shric/go-trpc/pkg/rm"
)

type dispatch map[string]func(*transmissionrpc.Client, []string)

type options struct {
	List list.Options `command:"list" alias:"l" description:"List torrents"`
	Add  add.Options  `command:"add" alias:"a" description:"Add torrents"`
	Rm   rm.Options   `command:"rm" alias:"r" description:"Remove torrents"`
}

func main() {

	var arguments = new(options)

	p := flags.NewParser(arguments, flags.Default)
	remaining, err := p.Parse()
	if err != nil {
		os.Exit(1)
	}
	start := time.Now()
	client := Connect()

	switch p.Active.Name {
	case "list":
		list.List(client, arguments.List, remaining)
	case "add":
		add.Add(client, arguments.Add, remaining)
	case "rm":
		rm.Rm(client, arguments.Rm, remaining)
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
