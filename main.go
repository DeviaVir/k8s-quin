package main

import (
	"log"
	"os"
	"time"

	"github.com/DeviaVir/k8s-quin/quin"
	"github.com/urfave/cli"
)

var (
	buildVersion     = "dev"
	buildSCMRevision = "0"
)

func main() {
	log.Println("starting quin")

	app := cli.NewApp()
	app.Name = "quin"
	app.Version = buildVersion + "-" + buildSCMRevision

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config-source",
			Value: "kubernetes",
			Usage: "Where to grab configuration from, valid options: kubernetes",
		},
		cli.BoolTFlag{
			Name:  "kubernetes-internal",
			Usage: "Use internal kubernetes configuration (set to false when you want to use the kubeconfig in your home dir).",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Print debug log lines.",
		},
		cli.IntFlag{
			Name:  "peer-poll-sec",
			Value: 60,
			Usage: "Time between updating k8s peers",
		},
		cli.IntFlag{
			Name:  "ping-frequency-sec",
			Value: 60,
			Usage: "Time between running ping probes across peers",
		},
		cli.IntFlag{
			Name:  "http-port",
			Value: 9666,
			Usage: "http port to listen on",
		},
	}

	app.Action = actionRun

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func actionRun(c *cli.Context) error {
	if c.Bool("debug") {
		log.Println("running with")
		log.Printf("--config-source: %v", c.String("config-source"))
		log.Printf("--ping-frequency-sec: %d", c.Int("ping-frequency-sec"))
		log.Printf("--kubernetes-internal: %t", c.BoolT("kubernetes-internal"))
		log.Printf("--http-port: %d", c.Int("http-port"))
		log.Printf("--debug: %t", c.Bool("debug"))
	}

	// One time register on startup. Can't record metrics until this completes.
	quin.RegisterMetrics(c.Bool("debug"))

	endpoints := quin.NewEndpointGrabber(c.String("config-source"), c.Int("peer-poll-sec"), c.BoolT("kubernetes-internal"))

	go func() {
		pingTicker := time.NewTicker(time.Duration(c.Int("ping-frequency-sec")) * time.Second)
		for {
			select {
			case _ = <-pingTicker.C:
				if c.Bool("debug") {
					log.Printf("Pinging peers")
				}
				quin.PingPeers(endpoints, c.Int("http-port"), c.Bool("debug"))
			}
		}
	}()

	quin.RunServer(c.Int("http-port"))

	return nil
}
