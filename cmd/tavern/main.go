package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
	"github.com/wybiral/tavern/app"
	"github.com/wybiral/torgo"
	"golang.org/x/crypto/ed25519"
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
		// onion commands
		cli.Command{
			Name:  "onion",
			Usage: "Tor onion tools",
			Subcommands: []cli.Command{
				cli.Command{
					Name:  "new",
					Usage: "Generate new Tor onion.key file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "type, t",
							Value: "ed25519",
							Usage: "Type of key to generate (rsa or ed25519)",
						},
					},
					Action: tavernOnionNew,
				},
			},
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
		os.Mkdir("public", os.ModePerm)
	}
	generateOnionFile("ed25519")
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

func tavernOnionNew(ctx *cli.Context) {
	generateOnionFile(ctx.String("type"))
}

func generateOnionFile(keyType string) {
	var onion *torgo.Onion
	_, err := os.Stat("onion.key")
	if err == nil {
		fmt.Printf("Overwrite existing onion.key file? (y/N): ")
		if !getConfirm() {
			return
		}
	}
	keyType = strings.ToLower(keyType)
	if keyType == "ed25519" {
		onion, err = generateEd25519()
		if err != nil {
			log.Fatal(err)
		}
	} else if keyType == "rsa" {
		onion, err = generateRSA()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("Invalid key type: " + keyType)
		os.Exit(1)
	}
	err = writeOnion(onion)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Generated: %s.onion\n", onion.ServiceID)
}

func getConfirm() bool {
	var s string
	_, err := fmt.Scan(&s)
	if err != nil {
		log.Fatal(err)
	}
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func generateEd25519() (*torgo.Onion, error) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	onion, err := torgo.OnionFromEd25519(key)
	if err != nil {
		return nil, err
	}
	return onion, nil
}

func generateRSA() (*torgo.Onion, error) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}
	onion, err := torgo.OnionFromRSA(key)
	if err != nil {
		return nil, err
	}
	return onion, nil
}

func writeOnion(onion *torgo.Onion) error {
	f, err := os.Create("onion.key")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(onion.PrivateKeyType + ":" + onion.PrivateKey)
	if err != nil {
		return err
	}
	return nil
}
