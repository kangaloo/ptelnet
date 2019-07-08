package main

import (
	_ "github.com/kangaloo/ptelnet/logger"
	"github.com/kangaloo/ptelnet/portscheck"
	"os"
)

// todo 增加summary功能

func main() {
	file := os.Args[1]
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	hosts, err := portscheck.NewHosts(f)
	if err != nil {
		panic(err)
	}

	hosts.Check()
}
