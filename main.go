package main

import (
	"os"

	"github.com/hekmon/transmissionrpc"
	"github.com/jessevdk/go-flags"
	cmd "github.com/shric/trpc/cmd/trpc"
)

type dispatch map[string]func(*transmissionrpc.Client, []string)

type options struct {
	List    cmd.ListOptions  `command:"list" alias:"l" description:"List torrents"`
	Add     cmd.AddOptions   `command:"add" alias:"a" description:"Add torrents"`
	Rm      cmd.RmOptions    `command:"rm" alias:"r" description:"Remove torrents"`
	Start   cmd.StartOptions `command:"start" description:"Start torrents"`
	Stop    cmd.StopOptions  `command:"stop" description:"Start torrents"`
	Version struct{}         `command:"version" description:"Print version"`
}

func main() {

	var arguments = new(options)

	p := flags.NewParser(arguments, flags.Default)
	remaining, err := p.Parse()
	if err != nil {
		os.Exit(1)
	}
	client := Connect()

	switch p.Active.Name {
	case "list":
		cmd.List(client, arguments.List, remaining)
	case "add":
		cmd.Add(client, arguments.Add, remaining)
	case "rm":
		cmd.Rm(client, arguments.Rm, remaining)
	case "stop":
		cmd.Stop(client, arguments.Stop, remaining)
	case "start":
		cmd.Start(client, arguments.Start, remaining)
	case "version":
		cmd.Version()
	}
}
