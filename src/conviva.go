//go:generate goversioninfo
package main

import (
	"fmt"
	"runtime"

	"os"

	sdk_args "github.com/newrelic/infra-integrations-sdk/v4/args"
	"github.com/newrelic/infra-integrations-sdk/v4/integration"
)

type argumentList struct {
	sdk_args.DefaultArgumentList
	ClientId		  string `help:"Conviva API client ID"`
	ClientSecret      string `help:"Conviva API client secret"`
	ConfigPath        string `help:"Path to YAML configuration"`
	ShowVersion       bool   `default:"false" help:"Print build information and exit"`
}

const (
	integrationName = "com.newrelic.odp.conviva"
)

var (
	args               argumentList
	integrationVersion = "0.0.0"
	gitCommit          = ""
	buildDate          = ""

)

func main() {	
	i, err := createIntegration()
	fatalIfErr(err)

	if args.ShowVersion {
		fmt.Printf(
			"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
			"Conviva",
			integrationVersion,
			fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
			runtime.Version(),
			gitCommit,
			buildDate)
		os.Exit(0)
	}

	log := i.Logger()

	log.Debugf(
		"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
		"Conviva",
		integrationVersion,
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		runtime.Version(),
		gitCommit,
		buildDate,
	)
	
	e, err := entity(i)
	fatalIfErr(err)

	/*
	if args.HasInventory() {
		fatalIfErr(setInventoryData(e.Inventory))
	}
	*/

	if args.ConfigPath == "" {
		fatalIfErr(fmt.Errorf("no config path specified"))
	}

	cfg, err := loadConfig(args.ConfigPath, log)
	fatalIfErr(err)

	if args.All() || args.HasMetrics() {
		log.Debugf("conviva metric collection enabled.")
		if len(cfg.Metrics) > 0 {
			initMetrics()
			err = getMetricsData(e, log, cfg)
			fatalIfErr(err)
		} else {
			log.Warnf("No metrics found to collect.")
		}
	}

	fatalIfErr(i.Publish())
}

func entity(i *integration.Integration) (*integration.Entity, error) {
	/*
	if args.RemoteMonitoring {
		hostname, port, err := parseStatusURL(args.StatusURL)
		if err != nil {
			return nil, err
		}
		n := fmt.Sprintf("%s:%s", hostname, port)
		return i.Entity(n, entityRemoteType)
	}
	*/
	return i.HostEntity, nil
}

func createIntegration() (*integration.Integration, error) {
	return integration.New(integrationName, integrationVersion, integration.Args(&args))
}