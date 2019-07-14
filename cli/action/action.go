package action

import (
	"errors"
	"github.com/kangaloo/ptelnet/portscheck"
	"github.com/kangaloo/ptelnet/util"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func Action(c *cli.Context) error {

	if c.IsSet("f") && c.IsSet("e") {
		return errors.New("can not use '-f' and '-e' flags at the same time")
	}

	if !c.IsSet("f") && !c.IsSet("e") {
		return cli.ShowAppHelp(c)
	}

	var reader io.ReadCloser

	if c.IsSet("f") {
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
		reader = f
	} else {
		r := strings.NewReader(c.String("e"))
		rc := ioutil.NopCloser(r)
		reader = rc
	}

	hosts, err := portscheck.NewHosts(reader)
	if err != nil {
		return err
	}

	// todo 将timeout参数放到Host结构里
	hosts.Check(c.Int("t"))
	return nil
}
