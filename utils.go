package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"hailo/conf"
	"hailo/hailo"
	"os"
	"time"
)

// printData prints the complete FDS data for the given endpoint
func printData(authServer string, username string, password string, fdsEndpoint string) {

	time.Sleep(time.Millisecond * 100) // Sleeps for a short time, since fmt output sometimes overlaps

	// Create config from command line args
	config := conf.BuildFdsConfig(authServer, username, password, fdsEndpoint)
	fmt.Println(" ---- Server config ----")
	pretty, _ := json.MarshalIndent(config, "", "\t")
	fmt.Println(string(pretty))

	// Get Specifications
	specs, _ := hailo.GetSpecs(config)
	for _, spec := range specs.Data {
		printSpec(config, spec)
		for _, subSpec := range spec.DeviceTypeSpecific.ComponentIdList {
			printSpec(config, subSpec)
		}
	}
}

// args holds command line arguments
type args struct {
	test        bool
	userName    string
	password    string
	authServer  string
	fdsEndpoint string
}

// determineArgs gets the
func determineArgs() args {
	var programArgs args
	flag.BoolVar(&programArgs.test, "t", false, "Test get data from Hailo FDS API without Eliona")
	flag.StringVar(&programArgs.userName, "user", "", "Username for Hailo FDS authentication endpoint (used for -t)")
	flag.StringVar(&programArgs.password, "password", "", "Password for Hailo FDS authentication endpoint (used for -t)")
	flag.StringVar(&programArgs.authServer, "auth", "", "Authentication endpoint for Hailo FDS (used for -t)")
	flag.StringVar(&programArgs.fdsEndpoint, "fds", "", "FDS endpoint for for Hailo FDS (used for -t)")

	flag.Usage = func() {
		fmt.Printf("Usage: \n")
		flag.PrintDefaults() // prints default usage
		os.Exit(0)
	}

	flag.Parse()
	return programArgs
}

// printSpec gets and prints out the data for the specification
func printSpec(config conf.Config, spec hailo.Spec) {
	fmt.Printf(" ---- Device %s ----\n", spec.DeviceId)
	pretty, _ := json.MarshalIndent(spec, "", "\t")
	fmt.Println(string(pretty))

	status, _ := hailo.GetStatus(config, spec.DeviceId)
	fmt.Printf(" ---- Status %s ----\n", spec.DeviceId)
	pretty, _ = json.MarshalIndent(status, "", "\t")
	fmt.Println(string(pretty))

	diag, _ := hailo.GetDiag(config, spec.DeviceId)
	fmt.Printf(" ---- Diagnostic %s ----\n", spec.DeviceId)
	pretty, _ = json.MarshalIndent(diag, "", "\t")
	fmt.Println(string(pretty))
}
