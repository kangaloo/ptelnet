package main

import (
	"github.com/kangaloo/ptelnet/portscheck"
	"os"
)

// todo 增加summary功能

func main() {
	file := "file.txt"
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
