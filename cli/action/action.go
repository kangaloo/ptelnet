package action

import (
	"errors"
	"github.com/kangaloo/ptelnet/portscheck"
	"github.com/kangaloo/ptelnet/util"
	"github.com/urfave/cli"
	"os"
)

func Action(c *cli.Context) error {
	if !c.IsSet("f") {
		return cli.ShowAppHelp(c)
	}

	file := c.String("f")

	if !util.IsExist(file) {
		return errors.New("file not exist")
	}

	if !util.IsFile(file) {
		return errors.New("not a file")
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	fi.IsDir()

	hosts, err := portscheck.NewHosts(f)
	if err != nil {
		return err
	}

	hosts.Check()
	return nil
}
