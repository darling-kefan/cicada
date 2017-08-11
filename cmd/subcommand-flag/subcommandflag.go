package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

func main() {
	globalFlags := []cli.Flag{
		cli.BoolFlag{Name: "debug, d", Usage: "Run in debug mode"},
	}

	adminServiceStatusCmd := cli.Command{
		Name: "status",
		//Flags: append([]cli.Flag{}, globalFlags...),
		Action: func(c *cli.Context) error {
			global := strconv.FormatBool(c.GlobalBool("debug"))
			local := strconv.FormatBool(c.Bool("debug"))
			fmt.Printf("%s: => local (%s), global (%s)\n", c.Command.Name, local, global)
			return nil
		},
	}

	adminServiceCmd := cli.Command{
		Name: "service",
		Flags: append([]cli.Flag{}, globalFlags...),
		Subcommands: []cli.Command{adminServiceStatusCmd},
	}

	adminCmd := cli.Command{
		Name: "admin",
		Flags: append([]cli.Flag{}, globalFlags...),
		Subcommands: []cli.Command{adminServiceCmd},
		Action: func(c *cli.Context) error {
			fmt.Printf("%s: => global (%#v)\n", c.Command.Name, c.GlobalBool("debug"))
			return nil
		},
	}

	app := cli.NewApp()
	app.Name = "lookup"
	app.Flags = append([]cli.Flag{}, globalFlags...)
	app.Commands = []cli.Command{adminCmd}

	app.Run(os.Args)
}
