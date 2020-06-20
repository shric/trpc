package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hekmon/transmissionrpc"
)

const (
	defaultPort    = 9091
	defaultTimeout = "30s"
)

func getHostPort() (string, uint16) {
	address, exists := os.LookupEnv("TR_HOST")
	host := "127.0.0.1"
	port := uint16(defaultPort)

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

func getAuth() (user string, pass string) {
	auth, exists := os.LookupEnv("TR_AUTH")
	if exists {
		x := strings.Split(auth, ":")
		user, pass = x[0], x[1]
	}

	return
}

// Connect returns a transmissionrpc.Client after connecting
// Parses environment variables:
//   TR_HOST: "host[:port]" (default port 9091)
//   TR_AUTH: "user:password"
func Connect() *transmissionrpc.Client {
	host, port := getHostPort()
	user, pass := getAuth()

	timeout, err := time.ParseDuration(defaultTimeout)
	if err != nil {
		fmt.Println(err)
	}

	transmissionbt, err := transmissionrpc.New(host, user, pass, &transmissionrpc.AdvancedConfig{
		HTTPS:       false,
		Port:        port,
		RPCURI:      "/transmission/rpc",
		HTTPTimeout: timeout,
		UserAgent:   "trpc"})
	if err != nil {
		panic(err)
	}

	ok, serverVersion, serverMinimumVersion, err := transmissionbt.RPCVersion()
	if err != nil {
		panic(err)
	}

	if !ok {
		panic(fmt.Sprintf("Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
	}

	return transmissionbt
}
