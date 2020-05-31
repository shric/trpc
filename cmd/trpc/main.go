package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shric/go-trpc/pkg/list"

	"github.com/hekmon/transmissionrpc"
)

func getHostPort() (string, uint16) {
	address, exists := os.LookupEnv("TR_HOST")
	host := "127.0.0.1"
	port := uint16(9091)
	if exists {
		x := strings.Split(address, ":")
		host = x[0]
		if len(x) > 1 {
			val, err := strconv.ParseInt(x[1], 10, 16)
			if err == nil {
				port = uint16(val)
			}
		}
	}
	return host, port
}

func getAuth() (string, string) {
	auth, exists := os.LookupEnv("TR_AUTH")
	if exists {
		x := strings.Split(auth, ":")
		return x[0], x[1]
	}
	return "", ""
}

func connect() *transmissionrpc.Client {
	host, port := getHostPort()
	user, pass := getAuth()
	timeout, err := time.ParseDuration("30s")
	if err != nil {
		fmt.Println(err)
	}
	transmissionbt, err := transmissionrpc.New(host, user, pass, &transmissionrpc.AdvancedConfig{
		HTTPS:       false,
		Port:        port,
		RPCURI:      "/transmission/rpc",
		HTTPTimeout: timeout,
		UserAgent:   "go-trpc"})
	if err != nil {
		fmt.Println(err)
	}
	return transmissionbt
}

func main() {
	start := time.Now()

	client := connect()

	list.List(client)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
