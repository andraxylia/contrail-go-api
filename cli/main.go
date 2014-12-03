//
// Copyright (c) 2014 Juniper Networks, Inc. All rights reserved.
//

package main

import (
	"github.com/Juniper/contrail-go-api"
	"flag"
	"fmt"
	"os"
)

type ExecFunc func(client *contrail.Client, flagSet *flag.FlagSet)

type CliCommand struct {
	flagSet *flag.FlagSet
	exec ExecFunc
}

var (
	/*
         * OpenContrail API server
         */
	oc_server string
	oc_port int

	/*
	 * Authentication
	 */
	os_auth_strategy string
	os_auth_url string
	os_tenant_name string
	os_tenant_id string
	os_username string
	os_password string
	os_region_name string
	os_token string

	commandMap map[string]CliCommand = make(map[string]CliCommand, 0)
)

func RegisterCliCommand(name string, flagSet *flag.FlagSet, exec ExecFunc) {
	commandMap[name] = CliCommand{flagSet, exec}
}

func InitFlags() {
	flag.StringVar(&oc_server, "server", "localhost",
		"OpenContrail API server hostname or address")
	flag.IntVar(&oc_port, "port", 8082,
		"OpenContrail API server port")

	default_strategy := os.Getenv("OS_AUTH_STRATEGY")
	if len(default_strategy) == 0 {
		default_strategy = "keystone"
	}
	flag.StringVar(&os_auth_strategy,
		"os-auth-strategy", default_strategy,
		"Authentication strategy (Env: OS_AUTH_STRATEGY)")
	flag.StringVar(&os_auth_url,
		"os-auth-url", os.Getenv("OS_AUTH_URL"),
		"Authentication URL (Env: OS_AUTH_URL)")
	flag.StringVar(&os_tenant_name,
		"os-tenant-name", os.Getenv("OS_TENANT_NAME"),
		"Authentication tenant name (Env: OS_TENANT_NAME)")
	flag.StringVar(&os_tenant_id,
		"os-tenant-id", os.Getenv("OS_TENANT_ID"),
		"Authentication tenant id (Env: OS_TENANT_ID)")
	flag.StringVar(&os_username,
		"os-username", os.Getenv("OS_USERNAME"),
		"Authentication username (Env: OS_USERNAME)")
	flag.StringVar(&os_password,
		"os-password", os.Getenv("OS_PASSWORD"),
		"Authentication password (Env: OS_PASSWORD)")
	flag.StringVar(&os_region_name,
		"os-region-name", os.Getenv("OS_REGION_NAME"),
		"Authentication region name (Env: OS_REGION_NAME)")
	flag.StringVar(&os_token,
		"os-token", os.Getenv("OS_TOKEN"),
		"Authentication URL (Env: OS_TOKEN)")
}

func usage() {
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "  Commands:\n")
	for name, _ := range commandMap {
		fmt.Fprintf(os.Stderr, "    %s\n", name)
	}
}

func main() {
	InitFlags()
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}

	command := flag.Arg(0)
	cmd, inMap := commandMap[command]
	if !inMap {
		usage()
		os.Exit(2)
	}

	flagSet := cmd.flagSet
	flagSet.Parse(flag.Args()[1:])

	client := contrail.NewClient(oc_server, oc_port)
	// TODO: setup authentication

	cmd.exec(client, flagSet)
}
