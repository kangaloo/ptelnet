package main

import (
	"github.com/kangaloo/ptelnet/cli/action"
	"github.com/kangaloo/ptelnet/cli/flag"
	_ "github.com/kangaloo/ptelnet/logger"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "ptelnet"
	app.Author = "Li Xiangyang"
	app.Email = "lixy4@belink.com"
	app.Usage = "A parallel telnet cli application"
	app.Flags = flag.GlobalFlags
	app.Action = action.Action

	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
