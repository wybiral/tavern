package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"github.com/wybiral/tavern/app"
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

func tavernHelp(ctx *cli.Context) {
	args := ctx.Args()
	if args.Present() {
		cli.ShowCommandHelp(ctx, args.First())
		return
	}
	cli.ShowAppHelp(ctx)
}

func tavernInit(ctx *cli.Context) {
	_, err := os.Stat("tavern.json")
	if err == nil {
		fmt.Println("Unable to init (tavern.json file already exists)")
		os.Exit(1)
	}
	c := app.DefaultConfig()
	err = c.WriteFile("tavern.json")
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Stat("public")
	if os.IsNotExist(err) {
		os.Mkdir("public", 0666)
	}
	fmt.Println("Done")
}

func tavernRun(ctx *cli.Context) {
	c := app.DefaultConfig()
	err := c.ReadFile("tavern.json")
	if os.IsNotExist(err) {
		fmt.Println("Missing tavern.json file")
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
	a, err := app.NewApp(c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Local server: http://%s:%d\n", c.Server.Host, c.Server.Port)
	if a.Tor != nil {
		fmt.Printf("Hidden service: http://%s.onion\n", a.Tor.Onion.ServiceID)
	}
	a.Run()
}
