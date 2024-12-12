package main

import (
	"fmt"
	"log"
	"os"
	"rde"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:   "Domain crawler",
		Usage:  "Crawl and enrich daily registred domains from ICAAN",
		Action: allCommand,
		Commands: []*cli.Command{
			{
				Name:   "load",
				Usage:  "Load data",
				Action: loadCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "output",
						Usage:    "Specify output folder. If not set, will be used OutputFolder from config.",
						Required: false,
					},
				},
			},
			{
				Name:   "regex",
				Usage:  "Search domains by regex",
				Action: regexCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Usage:    "Specify path to txt file with domains. If not set, will be used default filename after load command.",
						Required: false,
					},
				},
			},
			{
				Name:   "diff",
				Usage:  "Compare differences",
				Action: diffCommand,
			},
			{
				Name:   "all",
				Usage:  "Run all commands sequentially",
				Action: allCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func loadCommand(c *cli.Context) error {
	var (
		conf   *rde.Conf
		loader *rde.Loader
		err    error
	)

	conf, err = rde.NewConf()
	if err != nil {
		log.Fatal(err)
	}
	outputPath := c.String("output")
	if outputPath == "" {
		loader = rde.NewLoader(conf, nil)
	} else {
		loader = rde.NewLoader(conf, &outputPath)
	}

	st := time.Now()

	err = loader.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Finish loading for %d s", time.Since(st))

	return nil
}

func regexCommand(c *cli.Context) error {
	inputPath := c.String("input")
	if inputPath == "" {
		conf, err := rde.NewConf()
		if err != nil {
			panic(err)
		}
		inputPath = conf.Icaan.OutputFolder
	}
	fmt.Println("Running regex search on", inputPath)
	return nil
}

func diffCommand(c *cli.Context) error {
	fmt.Println("Running diff command...")
	return nil
}

func allCommand(c *cli.Context) error {
	fmt.Println("Running all commands...")
	return nil
}
