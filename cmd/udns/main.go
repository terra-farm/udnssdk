package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/logutils"

	"github.com/terra-farm/udnssdk"
)

var username string
var password string
var baseURL string
var logLevel string
var zone string
var domain string
var typ string

func init() {
	flag.StringVar(&username, "username", os.Getenv("ULTRADNS_USERNAME"), "ultradns username")
	flag.StringVar(&password, "password", os.Getenv("ULTRADNS_PASSWORD"), "ultradns password")
	flag.StringVar(&baseURL, "base-url", udnssdk.DefaultLiveBaseURL, "ultradns base url")
	flag.StringVar(&logLevel, "log-level", "WARN", "log level: DEBUG, WARN, ERROR. default: WARN")
	flag.StringVar(&zone, "zone", "", "dns zone")
	flag.StringVar(&domain, "domain", "", "dns domain")
	flag.StringVar(&typ, "type", "A", "dns type")
}

func main() {
	flag.Parse()

	if username == "" {
		fmt.Println("no username provided. Set with parameter -username=samdoe or environment variable ULTRADNS_USERNAME=samdoe")
		os.Exit(1)
	}
	if password == "" {
		fmt.Println("no password provided. Set with parameter -password=s3cr3t or environment variable ULTRADNS_PASSWORD=s3cr3t")
		os.Exit(1)
	}

	if zone == "" {
		fmt.Println("no zone provided. Set with parameter -zone=example.com.")
		os.Exit(1)
	}
	if domain == "" {
		fmt.Println("no domain provided. Set with parameter -domain=foo")
		os.Exit(1)
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)

	client, err := udnssdk.NewClient(username, password, baseURL)

	if err != nil {
		log.Fatalf("Error setting up client: %s", err)
	}

	k := udnssdk.RRSetKey{
		Zone: zone,
		Type: typ,
		Name: domain,
	}
	rrs, err := client.RRSets.Select(k)
	if err != nil {
		log.Fatalf("Error requesting records: %s", err)
	}
	data, err := json.MarshalIndent(rrs, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling response: %s", err)
	}
	fmt.Printf("%v\n", string(data))
}
