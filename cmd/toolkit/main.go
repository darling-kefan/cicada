package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/darling-kefan/cicada"
)

var adviewChannelCmd = cli.Command{
	Name: "adviewChannel",
	Usage: "The map of adview channels and tvm channels sync",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "127.0.0.1",
			Usage: "redis host",
		},
		cli.StringFlag{
			Name:  "port",
			Value: "6379",
			Usage: "redis port",
		},
	},
	Before: func(c *cli.Context) error {
		return nil
	},
	After: func(c *cli.Context) error {
		return nil
	},
	Action: cicada.AdviewChannelAction,
}

var detectCmd = cli.Command{
	Name: "detect",
	Usage: "Release detection tools",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Value: "ybt",
			Usage: "detect type, (ybt|tz|cp|ka)",
		},
		cli.StringFlag{
			Name: "source",
			Usage: "source(log) file",
		},
	},
	Before: func(c *cli.Context) error {
		return  nil
	},
	After: func(c *cli.Context) error {
		return nil
	},
	Action: cicada.DetectAction,
	OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "WRONG: %s\n", err.Error())
		return nil
	},
}

var clearCmd = cli.Command{
	Name: "clear",
	Usage: "clear tools",
	Subcommands: []cli.Command{
		cli.Command{
			Name: "redis6310",
			Usage: "clear redis keys, the port is 6310",
			Before: func(c *cli.Context) error {
				return  nil
			},
			After: func(c *cli.Context) error {
				return nil
			},
			Action: cicada.ClearRedis6310Action,
		},
	},
	OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "WRONG: %s\n", err.Error())
		return nil
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "toolkit"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "tangshouqiang",
			Email: "tangshouqiang@tvmining.com",
		},
	}
	app.Usage = "small tool kit or small script"
	app.Commands = []cli.Command{
		adviewChannelCmd,
		detectCmd,
		clearCmd,
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Thar be no command %q here.\n", command)
	}
	app.Run(os.Args)
}

