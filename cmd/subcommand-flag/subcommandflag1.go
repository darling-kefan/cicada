package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var debug = false

func setDebug(c *cli.Context) error {
	if c.IsSet("debug") {
		debug = true
	}
	return nil
}

func printDebug(c *cli.Context) error {	
	fmt.Printf("%s, %s, %s\n", debug, c.GlobalBool("debug"), c.Bool("debug"))
	return nil
}

func main() {
	globalFlags := []cli.Flag{
		cli.BoolFlag{Name: "debug, d", Usage: "Run in debug mode"},
	}

	adminServiceStatusCmd := cli.Command{
		Name: "status",
		Before: setDebug,
		Flags: append([]cli.Flag{}, globalFlags...),
		Action: printDebug,
	}

	adminServiceCmd := cli.Command{
		Name: "service",
		Before: setDebug,
		Flags: append([]cli.Flag{}, globalFlags...),
		Subcommands: []cli.Command{adminServiceStatusCmd},
		Action: printDebug,
	}

	adminCmd := cli.Command{
		Name: "admin",
		Before: setDebug,
		Flags: append([]cli.Flag{}, globalFlags...),
		Subcommands: []cli.Command{adminServiceCmd},
		Action: printDebug,
	}

	app := cli.NewApp()
	app.Name = "lookup"
	app.Flags = append([]cli.Flag{}, globalFlags...)
	app.Commands = []cli.Command{adminCmd}
	app.Before = setDebug

	app.Run(os.Args)
}
