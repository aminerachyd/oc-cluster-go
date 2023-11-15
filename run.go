package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	arg        string
	clusterUrl string
	username   string
	output     string
	delete     bool
}

func Run() {
	flags := parseFlags()
	if flags.arg == "" {
		fmt.Println("Should have at least one argument")
		return
	}

	config := ReadConfig()

	if flags.arg == "edit" {
		EditConfigInEditor()
	} else if flags.arg == "list" {
		ListClusters(flags.output, config)
	} else if flags.arg == "sync" {
		// TODO
	} else {
		if flags.clusterUrl != "" && flags.username != "" {
			AddCluster(flags.arg, flags.clusterUrl, flags.username, config)
			ConnectToCluster(flags.arg, config)
		} else if flags.delete {
			DeleteCluster(flags.arg, config)
		} else {
			ConnectToCluster(flags.arg, config)
		}
	}
}

func parseFlags() Flags {
	clusterUrl := flag.String("clusterUrl", "", "The URL of the cluster (API server)")
	username := flag.String("username", "", "The username to login on the cluster")
	output := flag.String("format", "wide", "The format in which to print the info")
	delete := flag.Bool("delete", false, "Delete the specified cluster")
	flag.Parse()
	clusterName := flag.Arg(0)

	return Flags{
		arg:        clusterName,
		clusterUrl: *clusterUrl,
		username:   *username,
		output:     *output,
		delete:     *delete,
	}
}
