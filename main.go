package main

import (
	cmd "github.com/shric/trpc/cmd/trpc"
)

func main() {
	client := Connect()
	cmd.Run(client)
}
