package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "gs",
		Usage:  "systemd manager",
		Action: default_launch,
		Commands: []*cli.Command{
			{
				Name:    "tui",
				Aliases: []string{"ui"},
				Usage:   "ui",
				Action:  main_tui,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func default_launch(c *cli.Context) error {
	fmt.Println("duang")
	return nil
}

func init(){
	file := "./" +"msg"+ ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}