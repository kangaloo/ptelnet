package flag

import "github.com/urfave/cli"

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "f, file",
		Usage: "specify the ip list `file`",
	},
}
