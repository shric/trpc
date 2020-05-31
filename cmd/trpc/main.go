package main

import (
	"fmt"
	"time"

	"github.com/shric/go-trpc/pkg/list"
)

func main() {
	start := time.Now()

	client := Connect()

	list.List(client)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
