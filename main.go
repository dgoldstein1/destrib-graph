package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"strconv"
)

// mock out logging calls for testing
var logFatalf = log.Fatalf
var logWarn = log.Warnf
var logMsg = log.Infof
var logErr = log.Errorf
var logDebug = log.Debugf

// checks environment for required env vars
func parseEnv() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	requiredEnvs := []string{
		"GRAPH_DB_STORE_PORT",
		"GRAPH_DOCS_DIR",
		"GRAPH_DB_ARANGO_ENDPOINT",
		"GRAPH_DB_NAME",
	}
	for _, v := range requiredEnvs {
		if os.Getenv(v) == "" {
			logFatalf("'%s' was not set", v)
		} else {
			// print out config
			logMsg("%s=%s", v, os.Getenv(v))
		}
	}
	i, err := strconv.Atoi(os.Getenv("GRAPH_DB_STORE_PORT"))
	if err != nil {
		logFatalf(err.Error())
	}
	if i < 1000 || i > 65535 {
		logFatalf("GRAPH_DB_STORE_PORT must be a valid port in range but was '%i'", i)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "destrib-graph"
	app.Usage = "highly-distributed directional graph"
	app.Description = "Fast and portable graph storage"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "servies graph database",
			Action: func(c *cli.Context) error {
				parseEnv()
				// port has already been validated
				r, _ := SetupRouter(os.Getenv("GRAPH_DOCS_DIR"))
				port := os.Getenv("GRAPH_DB_STORE_PORT")
				// run indefinitely
				return r.Run(fmt.Sprintf(":%s", port))
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
