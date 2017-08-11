package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/darling-kefan/cicada"
)

var uniformSpeedCmd = cli.Command{
	Name:  "uniformspeed",
	Usage: "Control ads with uniform show",
	Flags: []cli.Flag{
		cli.IntFlag{Name:"segment, s", Value:6, Usage:"How many segments generating in future"},
		cli.IntFlag{Name:"day, d", Value:0, Usage:"Generate one day of data, 0-tody, 1-tomorrow"},
	},
	Before: func(c *cli.Context) error {
		return nil
	},
	After: func(c *cli.Context) error {
		return nil
	},
	Action: cicada.UniformSpeedAction,
}

func main() {
	globalFlags := []cli.Flag{
		cli.StringFlag{Name:"env", Value:"dev", Usage:"Running Environment, {local|dev(develop)|prod(production)}"},
	}

	app := cli.NewApp()
	app.Name = "tvmads"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "tangshouqiang",
			Email: "tangshouqiang@tvmining.com",
		},
	}
	app.Usage = "Tvm ads data sync and operate"
	app.Flags = append([]cli.Flag{}, globalFlags...)
	app.Commands = []cli.Command{
		uniformSpeedCmd,
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Thar be no command %q here.\n", command)
	}
	app.Run(os.Args)
}
