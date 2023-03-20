package main

import (
	"github.com/newrelic/infra-integrations-sdk/v4/log"
)

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}