package flag

import "github.com/urfave/cli"

// todo 增加 -t 参数，设置超时时间
//  增加 -e 参数，测试单个ip:port，和 -f 参数冲突，输出详细的summary和result，不输出到文件

var GlobalFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "f, file",
		Usage: "specify the ip list `file`",
	},
	cli.IntFlag{
		Name:  "t, timeout",
		Usage: "specify the connect `timeout`",
		Value: 10,
	},
	cli.StringFlag{
		Name:  "e, endpoint",
		Usage: "test the specific `ip:port`",
	},
}
