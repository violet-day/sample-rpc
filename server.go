package main

import (
	"github.com/violet-day/simple-rpc/rpc"
)

func main() {
	s := &rpc.Server{}
	s.Listen(":8080")
}
