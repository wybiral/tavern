package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

// current tavern CLI version
const version = "0.0.1"

func main() {
	// hide default flags
	cli.HelpFlag = cli.StringFlag{Hidden: true}
	cli.VersionFlag = cli.StringFlag{Hidden: true}
	// setup CLI app
	c := cli.NewApp()
	c.CommandNotFound = func(ctx *cli.Context, command string) {
		fmt.Printf("Command not found: %v\n", command)
		os.Exit(1)
	}
	c.Version = version
	c.Usage = "tarven publishing system"
	c.Commands = []cli.Command{
		// help command
		cli.Command{
			Name:      "help",
			Usage:     "Shows all commands or help for one command",
			ArgsUsage: "[command]",
			Action:    tavernHelp,
		},
		// init command
		cli.Command{
			Name:   "init",
			Usage:  "Initialize tavern in current directory",
			Action: tavernInit,
		},
		// run command
		cli.Command{
			Name:   "run",
			Usage:  "Run tavern in current directory",
			Action: tavernRun,
		},
		// version command
		cli.Command{
			Name:  "version",
			Usage: "Print version",
			Action: func(ctx *cli.Context) {
				fmt.Println(c.Version)
			},
		},
	}
	// run CLI app
	err := c.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func tavernHelp(c *cli.Context) {
	args := c.Args()
	if args.Present() {
		cli.ShowCommandHelp(c, args.First())
		return
	}
	cli.ShowAppHelp(c)
}

func tavernInit(c *cli.Context) {
	log.Fatal("not implemented")
}

func tavernRun(c *cli.Context) {
	log.Fatal("not implemented")
}
